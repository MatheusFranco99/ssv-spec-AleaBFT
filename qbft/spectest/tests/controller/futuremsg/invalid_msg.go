package futuremsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidMsg tests future msg invalid msg
func InvalidMsg() *ControllerSyncSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     10,
		Round:      3,
		Identifier: identifier[:],
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})
	msg.Signature = nil

	return &ControllerSyncSpecTest{
		Name: "future invalid msg",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "3cf8cd7d050943781b25a611cc7b51bc051608f316556a843936b62c276984cb",
		ExpectedError:        "invalid future msg: invalid decided msg: message signature is invalid",
	}
}
