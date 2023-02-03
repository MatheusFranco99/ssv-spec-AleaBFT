package alea

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type ProposedValueCheckF func(data []byte) error
type ProposerF func(state *State, round Round) types.OperatorID

// Instance is a single QBFT instance that starts with a Start call (including a value).
// Every new msg the ProcessMsg function needs to be called
type Instance struct {
	State  *State
	config IConfig

	processMsgF *types.ThreadSafeF
	startOnce   sync.Once
	StartValue  []byte
	verbose     bool
	peers       map[types.OperatorID]*Instance
	msgMutex    sync.Mutex
}

func NewInstance(
	config IConfig,
	share *types.Share,
	identifier []byte,
	height Height,
) *Instance {
	return &Instance{
		State: &State{
			Share:             share,
			ID:                identifier,
			Round:             FirstRound,
			Height:            height,
			LastPreparedRound: NoRound,
			ProposeContainer:  NewMsgContainer(),
			BatchSize:         1,
			VCBCState:         NewVCBCState(),
			FillGapContainer:  NewMsgContainer(),
			FillerContainer:   NewMsgContainer(),
			AleaDefaultRound:  FirstRound,
			Delivered:         NewVCBCQueue(),
			StopAgreement:     false,
			ACState:           NewACState(),
			FillerMsgReceived: 0,
		},
		config:      config,
		processMsgF: types.NewThreadSafeF(),
		verbose:     false,
		peers:       make(map[types.OperatorID]*Instance),
		msgMutex:    sync.Mutex{},
	}
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height Height) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		go i.StartAgreementComponent()

		go i.ProposeCycle()
	})
}

func (i *Instance) ProposeCycle() {

	signedMessage, _ := CreateProposal(i.State, i.config, []byte{1, 2, 3, 4})

	time.Sleep(2 * time.Second)
	for {
		time.Sleep(time.Duration(uint64(i.State.Share.OperatorID)) * 3 * time.Second)
		i.ProcessMsg(signedMessage)
	}

}

func (i *Instance) RegisterPeer(inst *Instance, opID types.OperatorID) {
	i.peers[opID] = inst
}

func (i *Instance) Deliver(proposals []*ProposalData) int {

	for _, proposal := range proposals {
		unixTime := int64(binary.BigEndian.Uint64(proposal.Data))
		elapsed := (time.Now().UnixMicro() - unixTime)
		fmt.Println("VALUE:", unixTime, "; LATENCY in micro:", elapsed)
	}

	if len(i.State.Delivered.data) >= 20 {
		os.Exit(1)
	}

	return 1
}

func (i *Instance) Broadcast(msg *SignedMessage) error {
	for _, inst := range i.peers {
		go inst.ProcessMsg(msg)
	}
	return nil
}

func (i *Instance) SendTCP(msg *SignedMessage, operatorID types.OperatorID) error {
	if inst, exists := i.peers[operatorID]; exists {
		go inst.ProcessMsg(msg)
	}
	return nil
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {

	i.msgMutex.Lock()
	defer i.msgMutex.Unlock()

	if err := i.BaseMsgValidation(msg); err != nil {
		return false, nil, nil, errors.Wrap(err, "invalid signed message")
	}

	res := i.processMsgF.Run(func() interface{} {
		switch msg.Message.MsgType {
		case ProposalMsgType:
			return i.uponProposal(msg, i.State.ProposeContainer)
		case FillGapMsgType:
			return i.uponFillGap(msg, i.State.FillGapContainer)
		case FillerMsgType:
			return i.uponFiller(msg, i.State.FillerContainer)
		case ABAInitMsgType:
			return i.uponABAInit(msg)
		case ABAAuxMsgType:
			return i.uponABAAux(msg)
		case ABAConfMsgType:
			return i.uponABAConf(msg)
		case ABAFinishMsgType:
			return i.uponABAFinish(msg)
		case VCBCSendMsgType:
			return i.uponVCBCSend(msg)
		case VCBCReadyMsgType:
			return i.uponVCBCReady(msg)
		case VCBCFinalMsgType:
			return i.uponVCBCFinal(msg)
		case VCBCRequestMsgType:
			return i.uponVCBCRequest(msg)
		case VCBCAnswerMsgType:
			return i.uponVCBCAnswer(msg)
		default:
			return errors.New("signed message type not supported")
		}
	})
	if res != nil {
		return false, nil, nil, res.(error)
	}
	return i.State.Decided, i.State.DecidedValue, aggregatedCommit, nil
}

func (i *Instance) BaseMsgValidation(msg *SignedMessage) error {
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid signed message")
	}

	if msg.Message.Round < i.State.Round {
		return errors.New("past round")
	}

	switch msg.Message.MsgType {
	case ProposalMsgType:
		return isValidProposal(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case FillGapMsgType:
		return isValidFillGap(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case FillerMsgType:
		return isValidFiller(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case VCBCSendMsgType:
		return isValidVCBCSend(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case VCBCReadyMsgType:
		return isValidVCBCReady(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case VCBCFinalMsgType:
		return isValidVCBCFinal(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case VCBCRequestMsgType:
		return isValidVCBCRequest(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case VCBCAnswerMsgType:
		return isValidVCBCAnswer(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case ABAInitMsgType:
		return isValidABAInit(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case ABAAuxMsgType:
		return isValidABAAux(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case ABAConfMsgType:
		return isValidABAConf(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	case ABAFinishMsgType:
		return isValidABAFinish(
			i.State,
			i.config,
			msg,
			i.config.GetValueCheckF(),
			i.State.Share.Committee,
		)
	default:
		return errors.New("signed message type not supported")
	}
}

// IsDecided interface implementation
func (i *Instance) IsDecided() (bool, []byte) {
	if state := i.State; state != nil {
		return state.Decided, state.DecidedValue
	}
	return false, nil
}

// GetConfig returns the instance config
func (i *Instance) GetConfig() IConfig {
	return i.config
}

// GetHeight interface implementation
func (i *Instance) GetHeight() Height {
	return i.State.Height
}

// GetRoot returns the state's deterministic root
func (i *Instance) GetRoot() ([]byte, error) {
	return i.State.GetRoot()
}

// Encode implementation
func (i *Instance) Encode() ([]byte, error) {
	return json.Marshal(i)
}

// Decode implementation
func (i *Instance) Decode(data []byte) error {
	return json.Unmarshal(data, &i)
}
