package alea

import "sync"

type ACState struct {
	ACRound  ACRound
	ABAState map[ACRound]*ABAState
	mutex    sync.Mutex
}

func NewACState() *ACState {
	acState := &ACState{
		ACRound:  FirstACRound,
		ABAState: make(map[ACRound]*ABAState),
		mutex:    sync.Mutex{},
	}
	acState.ABAState[acState.ACRound] = NewABAState(acState.ACRound)
	return acState
}

func (s *ACState) IncrementRound() {
	s.mutex.Lock()
	// update info
	s.ACRound += 1
	s.mutex.Unlock()
	s.InitializeRound(s.ACRound)
}

func (s *ACState) InitializeRound(acRound ACRound) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.ABAState[acRound]; !exists {
		s.ABAState[acRound] = NewABAState(acRound)
	}
}

func (s *ACState) GetCurrentABAState() *ABAState {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.ABAState[s.ACRound]; !exists {
		s.ABAState[s.ACRound] = NewABAState(s.ACRound)
	}
	return s.ABAState[s.ACRound]
}

func (s *ACState) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ABAState[s.ACRound] = NewABAState(s.ACRound)
}

func (s *ACState) GetABAState(acRound ACRound) *ABAState {
	if _, exists := s.ABAState[acRound]; !exists {
		s.InitializeRound(acRound)
	}
	return s.ABAState[acRound]
}
