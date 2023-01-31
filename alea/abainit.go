package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAInit(signedABAInit *SignedMessage) error {

	if i.verbose {
		fmt.Println(Purple("uponABAInit"))
	}

	// get data
	abaInitData, err := signedABAInit.Message.GetABAInitData()
	if err != nil {
		return errors.Wrap(err, "uponABAInit: could not get abainitdata from signedABAInit")
	}

	// if future round -> intialize future state
	if abaInitData.ACRound > i.State.ACState.ACRound {
		i.State.ACState.InitializeRound(abaInitData.ACRound)
	}
	if abaInitData.Round > i.State.ACState.GetCurrentABAState().Round {
		i.State.ACState.GetCurrentABAState().InitializeRound(abaInitData.Round)
	}
	// old message -> ignore
	if abaInitData.ACRound < i.State.ACState.ACRound {
		return nil
	}
	if abaInitData.Round < i.State.ACState.GetCurrentABAState().Round {
		return nil
	}

	abaState := i.State.ACState.GetABAState(abaInitData.ACRound)

	// add the message to the container
	abaState.ABAInitContainer.AddMsg(signedABAInit)

	// sender
	senderID := signedABAInit.GetSigners()[0]

	alreadyReceived := abaState.hasInit(abaInitData.Round, senderID, abaInitData.Vote)
	if i.verbose {
		fmt.Println(Purple("\tsenderID:", senderID, ", vote:", abaInitData.Vote, ", round:", abaInitData.Round, ", already received before:", alreadyReceived))
	}
	// if never received this msg, update
	if !alreadyReceived {
		// Set received msg
		abaState.setInit(abaInitData.Round, senderID, abaInitData.Vote)

		if i.verbose {
			fmt.Println(Purple("\tupdated counter. Vote:", abaInitData.Vote, ". InitCounter:", abaState.InitCounter))
		}
	}

	// weak support -> send INIT
	// if never sent INIT(b) but reached PartialQuorum (i.e. f+1, weak support), send INIT(b)
	for _, vote := range []byte{0, 1} {
		if !abaState.sentInit(abaInitData.Round, vote) && abaState.countInit(abaInitData.Round, vote) >= i.State.Share.PartialQuorum {
			if i.verbose {
				fmt.Println(Purple("\tgot weak support for (and never sent):", vote))
			}
			// send INIT
			initMsg, err := CreateABAInit(i.State, i.config, vote, abaInitData.Round, abaInitData.ACRound)
			if err != nil {
				return errors.Wrap(err, "uponABAInit: failed to create ABA Init message after weak support")
			}
			if i.verbose {
				fmt.Println(Purple("\tsending INIT"))
			}
			i.Broadcast(initMsg)
			// update sent flag
			abaState.setSentInit(abaInitData.Round, vote, true)
			abaState.setInit(abaInitData.Round, i.State.Share.OperatorID, vote)
		}
	}

	// strong support -> send AUX
	// if never sent AUX(b) but reached Quorum (i.e. 2f+1, strong support), sends AUX(b) and add b to values
	for _, vote := range []byte{0, 1} {

		if !abaState.sentAux(abaInitData.Round, vote) && abaState.countInit(abaInitData.Round, vote) >= i.State.Share.Quorum {
			if i.verbose {
				fmt.Println(Purple("\tgot strong support and never sent AUX:", vote))
			}

			// append vote

			abaState.AddToValues(abaInitData.Round, vote)
			if i.verbose {
				fmt.Println(Purple("\tadded vote to local values for round", abaInitData.Round, ", values:", abaState.Values[abaInitData.Round]))
			}

			// sends AUX(b)
			auxMsg, err := CreateABAAux(i.State, i.config, vote, abaInitData.Round, abaInitData.ACRound)
			if err != nil {
				return errors.Wrap(err, "uponABAInit: failed to create ABA Aux message after strong init support")
			}
			if i.verbose {
				fmt.Println(Purple("\tsending ABAAux"))
			}
			i.Broadcast(auxMsg)

			// update sent flag
			abaState.setSentAux(abaInitData.Round, vote, true)
			abaState.setAux(abaInitData.Round, i.State.Share.OperatorID, vote)
		}
	}

	return nil
}

func isValidABAInit(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != ABAInitMsgType {
		return errors.New("msg type is not ABAInitMsgType")
	}
	if signedMsg.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedMsg.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}
	if err := signedMsg.Signature.VerifyByOperators(signedMsg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	ABAInitData, err := signedMsg.Message.GetABAInitData()
	if err != nil {
		return errors.Wrap(err, "could not get ABAInitData data")
	}
	if err := ABAInitData.Validate(); err != nil {
		return errors.Wrap(err, "ABAInitData invalid")
	}

	// vote
	vote := ABAInitData.Vote
	if vote != 0 && vote != 1 {
		return errors.New("vote different than 0 and 1")
	}

	return nil
}

func CreateABAInit(state *State, config IConfig, vote byte, round Round, acRound ACRound) (*SignedMessage, error) {
	ABAInitData := &ABAInitData{
		Vote:    vote,
		Round:   round,
		ACRound: acRound,
	}
	dataByts, err := ABAInitData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAInit: could not encode abainit data")
	}
	msg := &Message{
		MsgType:    ABAInitMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAInit: failed signing abainit msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
