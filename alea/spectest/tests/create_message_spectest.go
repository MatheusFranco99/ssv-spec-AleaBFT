package tests

import (
	"encoding/hex"
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/stretchr/testify/require"
)

const (
	CreateProposal      = "CreateProposal"
	CreateVCBC          = "CreateVCBC"
	CreateABA           = "CreateABA"
	CreateFillGap       = "CreateFillGap"
	CreateFiller        = "CreateFiller"
	CreateABAInit       = "CreateABAInit"
	CreateABAAux        = "CreateABAAux"
	CreateABAConf       = "CreateABAConf"
	CreateABAFinish     = "CreateABAFinish"
	CreateVCBCBroadcast = "CreateVCBCBroadcast"
	CreateVCBCSend      = "CreateVCBCSend"
	CreateVCBCReady     = "CreateVCBCReady"
	CreateVCBCFinal     = "CreateVCBCFinal"
	CreateVCBCRequest   = "CreateVCBCRequest"
	CreateVCBCAnswer    = "CreateVCBCAnswer"
)

type CreateMsgSpecTest struct {
	Name          string
	Value         []byte
	Proposals     []*alea.ProposalData
	Priority      alea.Priority
	Vote          byte
	Votes         []byte
	Author        types.OperatorID
	Entries       [][]*alea.ProposalData
	Priorities    []alea.Priority
	Round         alea.Round
	Hash          []byte
	Proof         types.Signature
	CreateType    string
	ExpectedRoot  string
	ExpectedError string
}

func (test *CreateMsgSpecTest) Run(t *testing.T) {
	var msg *alea.SignedMessage
	var lastErr error
	switch test.CreateType {
	case CreateProposal:
		msg, lastErr = test.createProposal()
	case CreateVCBC:
		msg, lastErr = test.createVCBC()
	case CreateABA:
		msg, lastErr = test.createABA()
	case CreateFillGap:
		msg, lastErr = test.createFillGap()
	case CreateFiller:
		msg, lastErr = test.createFiller()
	case CreateABAInit:
		msg, lastErr = test.createABAInit()
	case CreateABAAux:
		msg, lastErr = test.createABAAux()
	case CreateABAConf:
		msg, lastErr = test.createABAConf()
	case CreateABAFinish:
		msg, lastErr = test.createABAFinish()
	case CreateVCBCBroadcast:
		msg, lastErr = test.createVCBCBroadcast()
	case CreateVCBCSend:
		msg, lastErr = test.createVCBCSend()
	case CreateVCBCReady:
		msg, lastErr = test.createVCBCReady()
	case CreateVCBCFinal:
		msg, lastErr = test.createVCBCFinal()
	case CreateVCBCRequest:
		msg, lastErr = test.createVCBCRequest()
	case CreateVCBCAnswer:
		msg, lastErr = test.createVCBCAnswer()
	default:
		t.Fail()
	}

	r, err := msg.GetRoot()
	if err != nil {
		lastErr = err
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	require.EqualValues(t, test.ExpectedRoot, hex.EncodeToString(r))
}

func (test *CreateMsgSpecTest) createProposal() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateProposal(state, config, test.Value)
}

func (test *CreateMsgSpecTest) createVCBC() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBC(state, config, test.Proposals, test.Priority)
}

func (test *CreateMsgSpecTest) createABA() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateABA(state, config, test.Vote, test.Round)
}

func (test *CreateMsgSpecTest) createFillGap() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateFillGap(state, config, test.Author, test.Priority)
}

func (test *CreateMsgSpecTest) createFiller() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateFiller(state, config, test.Entries, test.Priorities, test.Author)
}

func (test *CreateMsgSpecTest) createABAInit() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateABAInit(state, config, test.Vote, test.Round)
}

func (test *CreateMsgSpecTest) createABAAux() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateABAAux(state, config, test.Vote, test.Round)
}

func (test *CreateMsgSpecTest) createABAConf() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateABAConf(state, config, test.Votes, test.Round)
}

func (test *CreateMsgSpecTest) createABAFinish() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateABAFinish(state, config, test.Vote)
}

func (test *CreateMsgSpecTest) createVCBCBroadcast() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBCBroadcast(state, config, test.Proposals, test.Priority, test.Author)
}

func (test *CreateMsgSpecTest) createVCBCSend() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBCSend(state, config, test.Proposals, test.Priority, test.Author)
}

func (test *CreateMsgSpecTest) createVCBCReady() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBCReady(state, config, test.Hash, test.Priority, test.Author)
}

func (test *CreateMsgSpecTest) createVCBCFinal() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBCFinal(state, config, test.Hash, test.Priority, test.Proof, test.Author)
}

func (test *CreateMsgSpecTest) createVCBCRequest() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBCRequest(state, config, test.Priority, test.Author)
}

func (test *CreateMsgSpecTest) createVCBCAnswer() (*alea.SignedMessage, error) {
	ks := testingutils.Testing4SharesSet()
	state := &alea.State{
		Share: testingutils.TestingShareAlea(ks),
		ID:    []byte{1, 2, 3, 4},
	}
	config := testingutils.TestingConfigAlea(ks)

	return alea.CreateVCBCAnswer(state, config, test.Proposals, test.Priority, test.Proof, test.Author)
}

func (test *CreateMsgSpecTest) TestName() string {
	return "alea create message " + test.Name
}
