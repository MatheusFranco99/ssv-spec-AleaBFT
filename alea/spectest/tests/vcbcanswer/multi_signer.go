package vcbcanswer

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiSigner tests a proposal msg with > 1 signers
func MultiSigner() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.MultiSignAleaMsg([]*bls.SecretKey{testingutils.Testing4SharesSet().Shares[1], testingutils.Testing4SharesSet().Shares[2]}, []types.OperatorID{1, 2}, &alea.Message{
			MsgType:    alea.VCBCAnswerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCAnswerDataBytes(tests.ProposalDataList, alea.FirstPriority, tests.AggregatedMsgBytes(types.OperatorID(1), alea.FirstPriority), types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "vcbcanswer multi signer",
		Pre:            pre,
		PostRoot:       "f831e111116961c37e9c383d1e6e3532e2ec3d1513dfe995f002f7be80e64a8c",
		InputMessages:  msgs,
		OutputMessages: []*alea.SignedMessage{},
		ExpectedError:  "invalid signed message: msg allows 1 signer",
		DontRunAC:      true,
	}
}
