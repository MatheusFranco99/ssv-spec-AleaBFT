package qbft

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// HistoricalInstanceCapacity represents the upper bound of InstanceContainer a controllerprocessmsg can process messages for as messages are not
// guaranteed to arrive in a timely fashion, we physically limit how far back the controllerprocessmsg will process messages for
const HistoricalInstanceCapacity int = 5

type InstanceContainer [HistoricalInstanceCapacity]*Instance

func (i InstanceContainer) FindInstance(height Height) *Instance {
	for _, inst := range i {
		if inst != nil {
			if inst.GetHeight() == height {
				return inst
			}
		}
	}
	return nil
}

// addNewInstance will add the new instance at index 0, pushing all other stored InstanceContainer one index up (ejecting last one if existing)
func (i *InstanceContainer) addNewInstance(instance *Instance) {
	for idx := HistoricalInstanceCapacity - 1; idx > 0; idx-- {
		i[idx] = i[idx-1]
	}
	i[0] = instance
}

// Controller is a QBFT coordinator responsible for starting and following the entire life cycle of multiple QBFT InstanceContainer
type Controller struct {
	Identifier []byte
	Height     Height // incremental Height for InstanceContainer
	// StoredInstances stores the last HistoricalInstanceCapacity in an array for message processing purposes.
	StoredInstances InstanceContainer
	// HigherReceivedMessages holds all msgs from a higher height
	HigherReceivedMessages *MsgContainer
	Domain                 types.DomainType
	Share                  *types.Share
	signer                 types.SSVSigner
	valueCheck             ProposedValueCheckF
	storage                Storage
	network                Network
	proposerF              ProposerF
}

func NewController(
	identifier []byte,
	share *types.Share,
	domain types.DomainType,
	signer types.SSVSigner,
	valueCheck ProposedValueCheckF,
	storage Storage,
	network Network,
	proposerF ProposerF,
) *Controller {
	return &Controller{
		Identifier:             identifier,
		Height:                 -1, // as we bump the height when starting the first instance
		Domain:                 domain,
		Share:                  share,
		StoredInstances:        InstanceContainer{},
		HigherReceivedMessages: NewMsgContainer(),
		signer:                 signer,
		valueCheck:             valueCheck,
		storage:                storage,
		network:                network,
		proposerF:              proposerF,
	}
}

// StartNewInstance will start a new QBFT instance, if can't will return error
func (c *Controller) StartNewInstance(value []byte) error {
	if err := c.canStartInstance(c.Height+1, value); err != nil {
		return errors.Wrap(err, "can't start new QBFT instance")
	}

	c.bumpHeight()
	newInstance := c.addAndStoreNewInstance()
	newInstance.Start(value, c.Height)

	return nil
}

// ProcessMsg processes a new msg, returns decided message or error
func (c *Controller) ProcessMsg(msg *SignedMessage) (*SignedMessage, error) {
	if err := c.baseMsgValidation(msg); err != nil {
		return nil, errors.Wrap(err, "invalid msg")
	}

	if msg.Message.Height > c.Height {
		return c.processFutureMsg(msg)
	} else {
		return c.processMsgCurrentInstance(msg)
	}
}

func (c *Controller) processMsgCurrentInstance(msg *SignedMessage) (*SignedMessage, error) {
	inst := c.InstanceForHeight(msg.Message.Height)
	if inst == nil {
		return nil, errors.New(fmt.Sprintf("instance not found"))
	}

	prevDecided, _ := inst.IsDecided()

	decided := false
	var decidedMsg *SignedMessage
	var err error
	if isDecidedMsg(c.Share, msg) {
		if err := validateDecided(
			inst.State.Height,
			c.GenerateConfig(),
			msg,
			c.Share.Committee,
		); err != nil {
			return nil, errors.Wrap(err, "invalid decided msg")
		}

		added, err := inst.State.CommitContainer.AddFirstMsgForSignerAndRound(msg)
		if inst == nil || !added {
			return nil, errors.New("could not add decided msg")
		}

		msgDecidedData, err := msg.Message.GetCommitData()
		if err != nil {
			return nil, errors.Wrap(err, "could not get msg decided data")
		}

		inst.State.Decided = true
		inst.State.DecidedValue = msgDecidedData.Data

		decided = true
		decidedMsg = msg
	} else {
		decided, _, decidedMsg, err = inst.ProcessMsg(msg)
		if err != nil {
			return nil, errors.Wrap(err, "could not process msg")
		}
	}

	// if previously Decided we do not return Decided true again
	if prevDecided {
		return nil, err
	}

	// save the highest Decided
	if !decided {
		return nil, nil
	}

	if err := c.saveAndBroadcastDecided(decidedMsg); err != nil {
		// TODO - we do not return error, should log?
	}
	return msg, nil
}

func (c *Controller) baseMsgValidation(msg *SignedMessage) error {
	// verify msg belongs to controllerprocessmsg
	if !bytes.Equal(c.Identifier, msg.Message.Identifier) {
		return errors.New(fmt.Sprintf("message doesn't belong to Identifier"))
	}

	return nil
}

func (c *Controller) InstanceForHeight(height Height) *Instance {
	return c.StoredInstances.FindInstance(height)
}

func (c *Controller) bumpHeight() {
	c.Height++
}

// GetIdentifier returns QBFT Identifier, used to identify messages
func (c *Controller) GetIdentifier() []byte {
	return c.Identifier
}

// addAndStoreNewInstance returns creates a new QBFT instance, stores it in an array and returns it
func (c *Controller) addAndStoreNewInstance() *Instance {
	i := NewInstance(c.GenerateConfig(), c.Share, c.Identifier, c.Height)
	c.StoredInstances.addNewInstance(i)
	return i
}

func (c *Controller) canStartInstance(height Height, value []byte) error {
	if height > FirstHeight {
		// check prev instance if prev instance is not the first instance
		inst := c.StoredInstances.FindInstance(height - 1)
		if inst == nil {
			return errors.New("could not find previous instance")
		}
		if decided, _ := inst.IsDecided(); !decided {
			return errors.New("previous instance hasn't Decided")
		}
	}

	// check value
	if err := c.valueCheck(value); err != nil {
		return errors.Wrap(err, "value invalid")
	}

	return nil
}

// GetRoot returns the state's deterministic root
func (c *Controller) GetRoot() ([]byte, error) {
	rootStruct := struct {
		Identifier             []byte
		Height                 Height
		InstanceRoots          [][]byte
		HigherReceivedMessages *MsgContainer
		Domain                 types.DomainType
		Share                  *types.Share
	}{
		Identifier:             c.Identifier,
		Height:                 c.Height,
		InstanceRoots:          make([][]byte, len(c.StoredInstances)),
		HigherReceivedMessages: c.HigherReceivedMessages,
		Domain:                 c.Domain,
		Share:                  c.Share,
	}

	for i, inst := range c.StoredInstances {
		if inst != nil {
			r, err := inst.GetRoot()
			if err != nil {
				return nil, errors.Wrap(err, "failed getting instance root")
			}
			rootStruct.InstanceRoots[i] = r
		}
	}

	marshaledRoot, err := json.Marshal(rootStruct)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret[:], nil
}

// Encode implementation
func (c *Controller) Encode() ([]byte, error) {
	return json.Marshal(c)
}

// Decode implementation
func (c *Controller) Decode(data []byte) error {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return errors.Wrap(err, "could not decode controllerprocessmsg")
	}

	config := c.GenerateConfig()
	for _, i := range c.StoredInstances {
		if i != nil {
			i.config = config
		}
	}
	return nil
}

func (c *Controller) saveAndBroadcastDecided(aggregatedCommit *SignedMessage) error {
	if err := c.storage.SaveHighestDecided(aggregatedCommit); err != nil {
		return errors.Wrap(err, "could not save decided")
	}

	// Broadcast Decided msg
	byts, err := aggregatedCommit.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided message")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   ControllerIdToMessageID(c.Identifier),
		Data:    byts,
	}
	if err := c.network.BroadcastDecided(msgToBroadcast); err != nil {
		// We do not return error here, just Log broadcasting error.
		return errors.Wrap(err, "could not broadcast decided")
	}
	return nil
}

func (c *Controller) GenerateConfig() IConfig {
	return &Config{
		Signer:      c.signer,
		SigningPK:   c.Share.ValidatorPubKey,
		Domain:      c.Domain,
		ValueCheckF: c.valueCheck,
		Storage:     c.storage,
		Network:     c.network,
		ProposerF:   c.proposerF,
	}
}
