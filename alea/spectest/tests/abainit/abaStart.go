package abainit

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// AbaStart tests if AC will send ABAInit0
func AbaStart() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(byte(1), alea.FirstRound, alea.FirstACRound),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "abainit abastart",
		Pre:           pre,
		PostRoot:      "ce4a9b59d8eb50cfb437ce09762faf19e083d304b6c664fdf4962f96ee80be20",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), alea.FirstRound, alea.FirstACRound),
			}),
		},
		DontRunAC: false,
	}
}
