package vcbcanswer

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// EmptyData tests the case in which the message has empty proposals
func EmptyData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.VCBCAnswerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCAnswerDataBytes([]*alea.ProposalData{}, alea.FirstPriority, tests.AggregatedMsgBytes(types.OperatorID(1), alea.FirstPriority), types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "vcbcanswer empty data",
		Pre:           pre,
		PostRoot:      "f831e111116961c37e9c383d1e6e3532e2ec3d1513dfe995f002f7be80e64a8c",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: VCBCAnswerData invalid: VCBCAnswerData: no proposals",
		DontRunAC:     true,
	}
}
