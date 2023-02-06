package launch

import (

	// spec "github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"

	// "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
)

func Launch(n int) *alea.Instance {
	inst := testingutils.BaseInstanceAleaN(types.OperatorID(n))
	inst.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)
	return inst
}
