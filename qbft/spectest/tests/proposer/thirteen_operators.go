package proposer

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// ThirteenOperators tests round-robin proposer selection for 13 member committee
func ThirteenOperators() *tests.RoundRobinSpecTest {
	var p types.OperatorID
	heights := make([]qbft.Height, 0)
	rounds := make([]qbft.Round, 0)
	proposers := make([]types.OperatorID, 0)
	for h := qbft.FirstHeight; h < 100; h++ {
		p = types.OperatorID(h%13) + 1
		for r := qbft.FirstRound; r < 100; r++ {
			heights = append(heights, h)
			rounds = append(rounds, r)
			proposers = append(proposers, p)

			p++
			if p == 14 {
				p = 1
			}
		}
	}

	//fmt.Printf("h:%v\nr:%v\np:%v\n", heights, rounds, proposers)

	return &tests.RoundRobinSpecTest{
		Name:      "13 member committee",
		Share:     testingutils.TestingShare(testingutils.Testing13SharesSet()),
		Heights:   heights,
		Rounds:    rounds,
		Proposers: proposers,
	}
}
