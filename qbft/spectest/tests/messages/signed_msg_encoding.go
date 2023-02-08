package messages

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SignedMessageEncoding tests encoding SignedMessage
func SignedMessageEncoding() *tests.MsgSpecTest {
	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data: testingutils.ProposalDataBytes(
			[]byte{1, 2, 3, 4},
			[]*qbft.SignedMessage{
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
			},
			[]*qbft.SignedMessage{
				testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
					MsgType:    qbft.RoundChangeMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
			},
		),
	})

	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "signed message encoding",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
	}
}
