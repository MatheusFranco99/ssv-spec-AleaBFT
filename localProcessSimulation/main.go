package main

import (
	"encoding/binary"
	"time"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	// spec "github.com/attestantio/go-eth2-client/spec/phase0"
	// "github.com/herumi/bls-eth-go-binary/bls"
)

func getTime() []byte {
	currentTime := time.Now().UnixNano()
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(currentTime))
	return timeBytes
}

func main() {

	inst1 := testingutils.BaseInstanceAleaN(types.OperatorID(1))
	inst2 := testingutils.BaseInstanceAleaN(types.OperatorID(2))
	inst3 := testingutils.BaseInstanceAleaN(types.OperatorID(3))
	inst4 := testingutils.BaseInstanceAleaN(types.OperatorID(4))

	inst1.RegisterPeer(inst2, types.OperatorID(2))
	inst1.RegisterPeer(inst3, types.OperatorID(3))
	inst1.RegisterPeer(inst4, types.OperatorID(4))

	inst2.RegisterPeer(inst1, types.OperatorID(1))
	inst2.RegisterPeer(inst3, types.OperatorID(3))
	inst2.RegisterPeer(inst4, types.OperatorID(4))

	inst3.RegisterPeer(inst1, types.OperatorID(1))
	inst3.RegisterPeer(inst2, types.OperatorID(2))
	inst3.RegisterPeer(inst4, types.OperatorID(4))

	inst4.RegisterPeer(inst1, types.OperatorID(1))
	inst4.RegisterPeer(inst3, types.OperatorID(3))
	inst4.RegisterPeer(inst2, types.OperatorID(2))

	inst1.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)
	inst2.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)
	inst3.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)
	inst3.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)

	for {

	}
}
