package filler

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// WrongData tests the case in which the message has proposals that mismatch with the aggreagtedMsgBytes's proposals
func WrongData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()

	msgs := []*alea.SignedMessage{
		testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
			MsgType:    alea.FillerMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.FillerDataBytes([][]*alea.ProposalData{tests.ProposalDataList2}, []alea.Priority{alea.FirstPriority}, [][]byte{tests.AggregatedMsgBytes(types.OperatorID(1), alea.FirstPriority)}, types.OperatorID(1)),
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "filler wrong data",
		Pre:           pre,
		PostRoot:      "84ecec5237cd4c1ca3ce3044e04e792a8abed2d470f29e1dd9416ac00511eec2",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: hash of proposals given doesn't match hash in the VCBCReadyData of the aggregated message",
		DontRunAC:     true,
	}
}
