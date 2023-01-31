package alea

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
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
			BatchSize:         2,
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
		verbose:     true,
	}
}

// Start is an interface implementation
func (i *Instance) Start(value []byte, height Height) {
	i.startOnce.Do(func() {
		i.StartValue = value
		i.State.Round = FirstRound
		i.State.Height = height

		go i.Listen()
		time.Sleep(2 * time.Second)
		// go i.StartAgreementComponent()

		// fmt.Println("Starting instance")

		// -> Init
		// state variables are initiated on constructor NewInstance (namely, queues and S)

		// -> Broadcast
		// The broadcast part runs as an instance receives proposal or vcbc messages
		// 		proposal message: is the message that a client sends to the node
		// 		vcbc message: is the broadcast a node does after receiving a batch size number of proposals

		// The agreement component consists of an infinite loop and we shall call it with another Thread
	})
}

func (i *Instance) Deliver(proposals []*ProposalData) int {

	for _, proposal := range proposals {
		unixTime := int64(binary.BigEndian.Uint64(proposal.Data))
		t := time.Unix(unixTime, 0)
		elapsed := time.Since(t)
		fmt.Println("$$$$$$$$$$$$$$$$$$GOT NEW LATENCY:", elapsed)
	}

	// FIX ME : to be implemented
	return 1
}

func (i *Instance) Listen() error {
	port := strconv.Itoa(8000 + int(i.State.Share.OperatorID))
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on localhost:" + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return err
		}
		go i.handleRequest(conn)
	}
}

func (i *Instance) handleRequest(conn net.Conn) {
	buf, read_err := ioutil.ReadAll(conn)
	if read_err != nil {
		fmt.Println("failed:", read_err)
		return
	}
	signedMessage := &SignedMessage{}
	signedMessage.Decode(buf)
	_, _, _, err := i.ProcessMsg(signedMessage)
	if err != nil {
		fmt.Println("Error processing msg:", err)
	}
}

func (i *Instance) Broadcast(msg *SignedMessage) error {
	byts, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode message")
	}

	// msgID := types.MessageID{}
	// copy(msgID[:], msg.Message.Identifier)

	// msgToBroadcast := &types.SSVMessage{
	// 	MsgType: types.SSVConsensusMsgType,
	// 	MsgID:   msgID,
	// 	Data:    byts,
	// }

	// if i.verbose {
	// 	fmt.Println("\tBroadcasting:", msg.Message.MsgType, msg.Message.Data)
	// }

	// return i.config.GetNetwork().Broadcast(msgToBroadcast)
	for _, operator := range i.State.Share.Committee {
		operatorID := operator.OperatorID
		if operatorID != i.State.Share.OperatorID {
			port := strconv.Itoa(8000 + int(operatorID))
			conn, err := net.Dial("tcp", "localhost:"+port)
			if err != nil {
				fmt.Println(err)
			} else {
				time.Sleep(1 / 20 * time.Second)
				_, write_err := conn.Write(byts)
				if write_err != nil {
					fmt.Println("failed:", write_err)
				}
				conn.(*net.TCPConn).CloseWrite()
			}
		}
	}
	return nil
}

func (i *Instance) SendTCP(msg *SignedMessage, operatorID types.OperatorID) error {
	byts, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode message")
	}

	// msgID := types.MessageID{}
	// copy(msgID[:], msg.Message.Identifier)

	// msgToBroadcast := &types.SSVMessage{
	// 	MsgType: types.SSVConsensusMsgType,
	// 	MsgID:   msgID,
	// 	Data:    byts,
	// }

	// if i.verbose {
	// 	fmt.Println("\tBroadcasting:", msg.Message.MsgType, msg.Message.Data)
	// }

	// return i.config.GetNetwork().Broadcast(msgToBroadcast)
	port := strconv.Itoa(8000 + int(operatorID))
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println(err)
	} else {
		_, write_err := conn.Write(byts)
		if write_err != nil {
			fmt.Println("failed:", write_err)
		}
		conn.(*net.TCPConn).CloseWrite()
	}
	return nil
}

// ProcessMsg processes a new QBFT msg, returns non nil error on msg processing error
func (i *Instance) ProcessMsg(msg *SignedMessage) (decided bool, decidedValue []byte, aggregatedCommit *SignedMessage, err error) {
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
