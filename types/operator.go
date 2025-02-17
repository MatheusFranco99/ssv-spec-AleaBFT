package types

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID uint64

type OperatorList []OperatorID

func (operators OperatorList) ToUint32List() []uint32 {
	l := make([]uint32, 0)
	for _, opID := range operators {
		l = append(l, uint32(opID))
	}
	return l
}

// Operator represents an SSV operator node
type Operator struct {
	OperatorID OperatorID
	PubKey     []byte
}

// GetPublicKey returns the public key with which the node is identified with
func (n *Operator) GetPublicKey() []byte {
	return n.PubKey
}

// GetID returns the node's ID
func (n *Operator) GetID() OperatorID {
	return n.OperatorID
}
