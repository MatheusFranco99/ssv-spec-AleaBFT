package testingutils

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type TestingNetwork struct {
	BroadcastedMsgs           []*types.SSVMessage
	SyncHighestDecidedCnt     int
	SyncHighestChangeRoundCnt int
	DecidedByRange            [2]alea.Height
}

func NewTestingNetwork() *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SSVMessage, 0),
	}
}

func (net *TestingNetwork) Broadcast(message *types.SSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}

func (net *TestingNetwork) SyncHighestDecided(identifier types.MessageID) error {
	net.SyncHighestDecidedCnt++
	return nil
}

//func (net *TestingNetwork) SyncHighestDecided() error {
//	return nil
//}

// SyncDecidedByRange will sync decided messages from-to (including them)
func (net *TestingNetwork) SyncDecidedByRange(identifier types.MessageID, from, to alea.Height) {
	net.DecidedByRange = [2]alea.Height{from, to}
}
