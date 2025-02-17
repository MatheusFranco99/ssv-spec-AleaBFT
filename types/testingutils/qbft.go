package testingutils

import (
	"bytes"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

var TestingConfig = func(keySet *TestKeySet) *alea.Config {
	return &alea.Config{
		Signer:    NewTestingKeyManager(),
		SigningPK: keySet.Shares[1].GetPublicKey().Serialize(),
		Domain:    types.PrimusTestnet,
		ValueCheckF: func(data []byte) error {
			if bytes.Equal(data, TestingInvalidValueCheck) {
				return errors.New("invalid value")
			}

			// as a base validation we do not accept nil values
			if len(data) == 0 {
				return errors.New("invalid value")
			}
			return nil
		},
		ProposerF: func(state *alea.State, round alea.Round) types.OperatorID {
			return 1
		},
		Network: NewTestingNetwork(),
		Timer:   NewTestingTimer(),
	}
}

var TestingInvalidValueCheck = []byte{1, 1, 1, 1}

var TestingShare = func(keysSet *TestKeySet) *types.Share {
	return &types.Share{
		OperatorID:          1,
		ValidatorPubKey:     keysSet.ValidatorPK.Serialize(),
		SharePubKey:         keysSet.Shares[1].GetPublicKey().Serialize(),
		DomainType:          types.PrimusTestnet,
		Quorum:              keysSet.Threshold,
		PartialQuorum:       keysSet.PartialThreshold,
		Committee:           keysSet.Committee(),
		FeeRecipientAddress: TestingFeeRecipient,
	}
}

var BaseInstance = func() *alea.Instance {
	return baseInstance(TestingShare(Testing4SharesSet()), Testing4SharesSet(), []byte{1, 2, 3, 4})
}

var SevenOperatorsInstance = func() *alea.Instance {
	return baseInstance(TestingShare(Testing7SharesSet()), Testing7SharesSet(), []byte{1, 2, 3, 4})
}

var TenOperatorsInstance = func() *alea.Instance {
	return baseInstance(TestingShare(Testing10SharesSet()), Testing10SharesSet(), []byte{1, 2, 3, 4})
}

var ThirteenOperatorsInstance = func() *alea.Instance {
	return baseInstance(TestingShare(Testing13SharesSet()), Testing13SharesSet(), []byte{1, 2, 3, 4})
}

var baseInstance = func(share *types.Share, keySet *TestKeySet, identifier []byte) *alea.Instance {
	ret := alea.NewInstance(TestingConfig(keySet), share, identifier, alea.FirstHeight)
	ret.StartValue = []byte{1, 2, 3, 4}
	return ret
}

func NewTestingQBFTController(
	identifier []byte,
	share *types.Share,
	config qbft.IConfig,
) *qbft.Controller {
	return qbft.NewController(
		identifier,
		share,
		types.PrimusTestnet,
		config,
	)
}
