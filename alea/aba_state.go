package alea

import (
	"fmt"
	"sync"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

type ABAState struct {
	// message containers
	ABAInitContainer   *MsgContainer
	ABAAuxContainer    *MsgContainer
	ABAConfContainer   *MsgContainer
	ABAFinishContainer *MsgContainer
	// message counters
	InitCounter   map[Round]map[byte][]types.OperatorID
	AuxCounter    map[Round]map[byte][]types.OperatorID
	ConfCounter   map[Round][]types.OperatorID
	FinishCounter map[byte][]types.OperatorID
	// already sent message flags
	SentInit   map[Round][]bool
	SentAux    map[Round][]bool
	SentConf   map[Round]bool
	SentFinish []bool
	// current ABA round
	ACRound ACRound
	// value inputed to ABA
	Vin map[Round]byte
	// value decided by ABA
	Vdecided byte
	// current ABA round
	Round Round
	// values that completed strong support of INIT messages
	Values map[Round][]byte
	// terminate channel to announce to ABA caller
	Terminate bool
	mutex     sync.Mutex
}

func NewABAState(acRound ACRound) *ABAState {
	abaState := &ABAState{
		ABAInitContainer:   NewMsgContainer(),
		ABAAuxContainer:    NewMsgContainer(),
		ABAConfContainer:   NewMsgContainer(),
		ABAFinishContainer: NewMsgContainer(),
		InitCounter:        make(map[Round]map[byte][]types.OperatorID),
		AuxCounter:         make(map[Round]map[byte][]types.OperatorID),
		ConfCounter:        make(map[Round][]types.OperatorID),
		FinishCounter:      make(map[byte][]types.OperatorID),
		SentInit:           make(map[Round][]bool),
		SentAux:            make(map[Round][]bool),
		SentConf:           make(map[Round]bool),
		SentFinish:         make([]bool, 2),
		ACRound:            acRound,
		Vin:                make(map[Round]byte),
		Vdecided:           byte(2),
		Round:              FirstRound,
		Values:             make(map[Round][]byte),
		mutex:              sync.Mutex{},
	}

	abaState.InitializeRound(FirstRound)
	abaState.FinishCounter[0] = make([]types.OperatorID, 0)
	abaState.FinishCounter[1] = make([]types.OperatorID, 0)

	return abaState
}

func (s *ABAState) String() string {
	return fmt.Sprintf("ABAState{Round:%d, InitCounter:%v, AuxCounter:%v, ConfCounter:%v, FinishCounter:%v, SentInit:%v, SentAux:%v, SentConf:%v, SentFinish:%v, ACRound:%d, Vin:%v, values:%v}", s.Round, s.InitCounter, s.AuxCounter, s.ConfCounter, s.FinishCounter, s.SentInit, s.SentAux, s.SentConf, s.SentFinish, s.ACRound, s.Vin, s.Values)
}

func (s *ABAState) Coin(round Round) byte {
	// FIX ME : implement a RANDOM coin generator given the round number
	return byte(round % 2)
}

func (s *ABAState) InitializeRound(round Round) {

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.InitCounter[round]; !exists {
		s.InitCounter[round] = make(map[byte][]types.OperatorID)
		s.InitCounter[round][0] = make([]types.OperatorID, 0)
		s.InitCounter[round][1] = make([]types.OperatorID, 0)
	}

	if _, exists := s.AuxCounter[round]; !exists {
		s.AuxCounter[round] = make(map[byte][]types.OperatorID, 2)
		s.AuxCounter[round][0] = make([]types.OperatorID, 0)
		s.AuxCounter[round][1] = make([]types.OperatorID, 0)
	}

	if _, exists := s.ConfCounter[round]; !exists {
		s.ConfCounter[round] = make([]types.OperatorID, 0)
	}

	if _, exists := s.SentInit[round]; !exists {
		s.SentInit[round] = make([]bool, 2)
	}
	if _, exists := s.SentAux[round]; !exists {
		s.SentAux[round] = make([]bool, 2)
	}
	if _, exists := s.SentConf[round]; !exists {
		s.SentConf[round] = false
	}

	if _, exists := s.Values[round]; !exists {
		s.Values[round] = make([]byte, 0)
	}
}

func (s *ABAState) IncrementRound() {
	s.mutex.Lock()
	// update info
	s.Round += 1
	s.mutex.Unlock()
	s.InitializeRound(s.Round)
}

func (s *ABAState) hasInit(round Round, operatorID types.OperatorID, vote byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, opID := range s.InitCounter[round][vote] {
		if opID == operatorID {
			return true
		}
	}
	return false
}
func (s *ABAState) hasAux(round Round, operatorID types.OperatorID, vote byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, opID := range s.AuxCounter[round][vote] {
		if opID == operatorID {
			return true
		}
	}
	return false
}
func (s *ABAState) hasConf(round Round, operatorID types.OperatorID) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, opID := range s.ConfCounter[round] {
		if opID == operatorID {
			return true
		}
	}
	return false
}
func (s *ABAState) hasFinish(operatorID types.OperatorID) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, vote := range []byte{0, 1} {
		for _, opID := range s.FinishCounter[vote] {
			if opID == operatorID {
				return true
			}
		}
	}
	return false
}

func (s *ABAState) countInit(round Round, vote byte) uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return uint64(len(s.InitCounter[round][vote]))
}
func (s *ABAState) countAux(round Round, vote byte) uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return uint64(len(s.AuxCounter[round][vote]))
}
func (s *ABAState) countConf(round Round) uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return uint64(len(s.ConfCounter[round]))
}
func (s *ABAState) countFinish(vote byte) uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return uint64(len(s.FinishCounter[vote]))
}

func (s *ABAState) setInit(round Round, operatorID types.OperatorID, vote byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.InitCounter[round][vote]; !exists {
		s.InitCounter[round][vote] = make([]types.OperatorID, 0)
	}

	s.InitCounter[round][vote] = append(s.InitCounter[round][vote], operatorID)
}
func (s *ABAState) setAux(round Round, operatorID types.OperatorID, vote byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.AuxCounter[round][vote]; !exists {
		s.AuxCounter[round][vote] = make([]types.OperatorID, 0)
	}

	s.AuxCounter[round][vote] = append(s.AuxCounter[round][vote], operatorID)
}
func (s *ABAState) setConf(round Round, operatorID types.OperatorID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.ConfCounter[round]; !exists {
		s.ConfCounter[round] = make([]types.OperatorID, 0)
	}

	s.ConfCounter[round] = append(s.ConfCounter[round], operatorID)
}
func (s *ABAState) setFinish(operatorID types.OperatorID, vote byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.FinishCounter[vote]; !exists {
		s.FinishCounter[vote] = make([]types.OperatorID, 0)
	}

	s.FinishCounter[vote] = append(s.FinishCounter[vote], operatorID)
}

func (s *ABAState) sentInit(round Round, vote byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.SentInit[round][vote]
}
func (s *ABAState) sentAux(round Round, vote byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.SentAux[round][vote]
}
func (s *ABAState) sentConf(round Round) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.SentConf[round]
}
func (s *ABAState) sentFinish(vote byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.SentFinish[vote]
}

func (s *ABAState) setSentInit(round Round, vote byte, value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SentInit[round][vote] = value
}
func (s *ABAState) setSentAux(round Round, vote byte, value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SentAux[round][vote] = value
}
func (s *ABAState) setSentConf(round Round, value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SentConf[round] = value
}
func (s *ABAState) setSentFinish(vote byte, value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SentFinish[vote] = value
}

func (s *ABAState) GetValues(round Round) []byte {
	return s.Values[round]
}

func (s *ABAState) AddToValues(round Round, vote byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, value := range s.Values[round] {
		if value == vote {
			return
		}
	}
	s.Values[round] = append(s.Values[round], vote)
}

func (s *ABAState) isContainedInValues(round Round, values []byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	num_equal := 0
	for _, value := range values {
		for _, storedValue := range s.Values[round] {
			if value == storedValue {
				num_equal += 1
			}
		}
	}
	return num_equal == len(values)
}
func (s *ABAState) existsInValues(round Round, value byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, storedValue := range s.Values[round] {
		if value == storedValue {
			return true
		}
	}
	return false
}

func (s *ABAState) setVInput(round Round, vote byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Vin[round] = vote
}
func (s *ABAState) getVInput(round Round) byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Vin[round]
}

func (s *ABAState) setDecided(vote byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Vdecided = vote
}

func (s *ABAState) setTerminate(value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Terminate = value
}
