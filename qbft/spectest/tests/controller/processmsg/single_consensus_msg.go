package processmsg

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SingleConsensusMsg tests process msg of a single msg
func SingleConsensusMsg() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	return &tests.ControllerSpecTest{
		Name: "single consensus msg",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], 1, &qbft.Message{
						MsgType:    qbft.ProposalMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
					}),
				},
				ControllerPostRoot: "ede93c772804aff585990dc09b29d841ce55e024831ea05674960b126922a0af",
			},
		},
	}
}
