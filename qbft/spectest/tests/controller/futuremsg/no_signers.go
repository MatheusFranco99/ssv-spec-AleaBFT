package futuremsg

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// NoSigners tests future msg with no signers
func NoSigners() *ControllerSyncSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.Shares[3], 3, &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     10,
		Round:      3,
		Identifier: identifier[:],
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})
	msg.Signers = []types.OperatorID{}

	return &ControllerSyncSpecTest{
		Name: "future msgs no signer",
		InputMessages: []*qbft.SignedMessage{
			msg,
		},
		SyncDecidedCalledCnt: 0,
		ControllerPostRoot:   "6bd17213f8e308190c4ebe49a22ec00c91ffd4c91a5515583391e9977423370f",
		ExpectedError:        "invalid future msg: invalid decided msg: message signers is empty",
	}
}
