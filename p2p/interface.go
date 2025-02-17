package p2p

import "github.com/MatheusFranco99/ssv-spec-AleaBFT/types"

// Broadcaster is the interface used to abstract message broadcasting
type Broadcaster interface {
	Broadcast(message *types.SSVMessage) error
}

// Subscriber is used to abstract topic management
type Subscriber interface {
	Subscribe(vpk types.ValidatorPK) error
}
