package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// JustificationInvalidSig tests a single prepare justification msg with wrong signature
// TODO<olegshmuelov>: inputroot related, verify signature for justification
func JustificationInvalidSig() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	signQBFTMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg2 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	signQBFTMsg3 := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(2), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	rcMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height:        qbft.FirstHeight,
		Round:         2,
		Input:         []byte{1, 2, 3, 4},
		PreparedRound: qbft.FirstRound,
	})

	prepareMsgHeader, _ := signQBFTMsg.ToSignedMessageHeader()
	prepareMsgHeader2, _ := signQBFTMsg2.ToSignedMessageHeader()
	prepareMsgHeader3, _ := signQBFTMsg3.ToSignedMessageHeader()

	prepareJustifications := []*qbft.SignedMessageHeader{
		prepareMsgHeader,
		prepareMsgHeader2,
		prepareMsgHeader3,
	}
	rcMsg.RoundChangeJustifications = prepareJustifications
	rcMsgEncoded, _ := rcMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusRoundChangeMsgType),
			Data: rcMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "justification invalid sig",
		Pre:              pre,
		PostRoot:         "a8b80879ebf2ecee42fddc69b67dd5f6adfd6aa8b7114246aec80ce1bfef513a",
		InputMessagesSIP: msgs,
		OutputMessages:   []*qbft.SignedMessage{},
		ExpectedError:    "round change msg invalid: round change justification invalid: prepare msg signature invalid: failed to verify signature",
	}
}
