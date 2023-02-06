package abafinish

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Duplicate tests an ABAFinish repeated
func Duplicate() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{}
	signedMsg := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(2)], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(byte(1), alea.FirstACRound),
	})
	msgs = append(msgs, signedMsg)
	signedMsg = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(2)], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(byte(1), alea.FirstACRound),
	})
	msgs = append(msgs, signedMsg)
	signedMsg = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(2)], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(byte(1), alea.FirstACRound),
	})
	msgs = append(msgs, signedMsg)

	return &tests.MsgProcessingSpecTest{
		Name:          "abafinish duplicate",
		Pre:           pre,
		PostRoot:      "89559c73eed75e25ebfafaca0481dfa8453646fa30c91c2f1990214ce3a9b64c",
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
