package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CreateProposalPreviouslyPrepared tests creating a proposal msg,previously prepared
func CreateProposalPreviouslyPrepared() *tests.CreateMsgSpecTest {
	return &tests.CreateMsgSpecTest{
		CreateType: tests.CreateProposal,
		Name:       "create proposal previously prepared",
		Value:      &qbft.Data{Root: [32]byte{1, 2, 3, 4}, Source: []byte{1, 2, 3, 4}},
		RoundChangeJustifications: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  2,
				Input:  &qbft.Data{},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  2,
				Input:  &qbft.Data{},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  2,
				Input:  &qbft.Data{},
			}),
		},
		PrepareJustifications: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}},
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
				Height: qbft.FirstHeight,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: [32]byte{1, 2, 3, 4}},
			}),
		},
		ExpectedRoot: "0102030400000000000000000000000000000000000000000000000000000000",
	}
}
