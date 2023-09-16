package testingutils

import "github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"

func UnknownDutyValueCheck() alea.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}
