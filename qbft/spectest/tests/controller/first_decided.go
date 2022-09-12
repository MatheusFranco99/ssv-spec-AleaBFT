package controller

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FirstDecided tests a process msg after which first time decided
func FirstDecided() *tests.ControllerSpecTest {
	//identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	identifier := types.NewBaseMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "first decided",
		RunInstanceData: []struct {
			InputValue    []byte
			InputMessages []*types.Message
			Decided       bool
			DecidedVal    []byte
			DecidedCnt    uint
			SavedDecided  *qbft.SignedMessage
		}{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier, qbft.FirstHeight, testingutils.Testing4SharesSet()),
				Decided:       true,
				DecidedVal:    []byte{1, 2, 3, 4},
				DecidedCnt:    1,
			},
		},
	}
}
