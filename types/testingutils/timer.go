package testingutils

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
)

type TimerState struct {
	Timeouts int
	Round    alea.Round
}

type TestQBFTTimer struct {
	State TimerState
}

func NewTestingTimer() alea.Timer {
	return &TestQBFTTimer{
		State: TimerState{},
	}
}

func (t *TestQBFTTimer) TimeoutForRound(round alea.Round) {
	t.State.Timeouts++
	t.State.Round = round
}
