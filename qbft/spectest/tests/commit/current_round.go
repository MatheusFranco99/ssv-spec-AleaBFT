package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CurrentRound tests a commit msg with current round, should process
func CurrentRound() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	signMsg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		Height: qbft.FirstHeight,
		Round:  qbft.FirstRound,
		Input:  []byte{1, 2, 3, 4},
	})
	pre.State.ProposalAcceptedForCurrentRound = signMsg
	signMsgEncoded, _ := signMsg.Encode()

	msgs := []*types.Message{
		{
			ID:   types.PopulateMsgType(pre.State.ID, types.ConsensusCommitMsgType),
			Data: signMsgEncoded,
		},
	}

	return &tests.MsgProcessingSpecTest{
		Name:             "commit current round",
		Pre:              pre,
		PostRoot:         "952dc40814611c59abfa2cbc0445c30cc54646da48ccdac01d8b48943770c569",
		InputMessagesSIP: msgs,
	}
}
