package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAConf(signedABAConf *SignedMessage) error {
	if i.verbose {
		fmt.Println("uponABAConf")
	}
	// get data
	ABAConfData, err := signedABAConf.Message.GetABAConfData()
	if err != nil {
		return errors.Wrap(err, "uponABAConf:could not get ABAConfData from signedABAConf")
	}

	// if future round -> intialize future state
	if ABAConfData.ACRound > i.State.ACState.ACRound {
		i.State.ACState.InitializeRound(ABAConfData.ACRound)
	}
	if ABAConfData.Round > i.State.ACState.GetCurrentABAState().Round {
		i.State.ACState.GetCurrentABAState().InitializeRound(ABAConfData.Round)
	}
	// old message -> ignore
	if ABAConfData.ACRound < i.State.ACState.ACRound {
		return nil
	}
	if ABAConfData.Round < i.State.ACState.GetCurrentABAState().Round {
		return nil
	}

	abaState := i.State.ACState.GetABAState(ABAConfData.ACRound)

	// add the message to the containers
	abaState.ABAConfContainer.AddMsg(signedABAConf)

	// sender
	senderID := signedABAConf.GetSigners()[0]

	alreadyReceived := abaState.hasConf(ABAConfData.Round, senderID)
	if i.verbose {
		fmt.Println("\tsenderID:", senderID, ", votes:", ABAConfData.Votes, ", round:", ABAConfData.Round, ", already received before:", alreadyReceived)
	}
	// if never received this msg, update
	if !alreadyReceived {

		// determine if votes list is contained in local round values list
		isContained := abaState.isContainedInValues(ABAConfData.Round, ABAConfData.Votes)
		// list is contained -> update CONF counter
		if isContained {
			abaState.setConf(ABAConfData.Round, senderID)
			if i.verbose {
				fmt.Println("\tupdated confcounter:", abaState.countConf(ABAConfData.Round))
			}
		}
	}

	// reached strong support -> try to decide value
	if abaState.countConf(ABAConfData.Round) >= i.State.Share.Quorum {
		if i.verbose {
			fmt.Println("\treached quorum")
		}

		// get common coin
		s := abaState.Coin(abaState.Round)
		if i.verbose {
			fmt.Println("\tcoin:", s)
		}

		// if values = {0,1}, choose randomly (i.e. coin) value for next round
		if len(abaState.Values[ABAConfData.Round]) == 2 {

			abaState.setVInput(ABAConfData.Round+1, s)
			if i.verbose {
				fmt.Println("\tlength of values is 2", abaState.Values[ABAConfData.Round], "-> storing coin to next Vin")
			}
		} else {
			if i.verbose {
				fmt.Println("\tlength of values is 1:", abaState.Values[ABAConfData.Round])
			}
			abaState.setVInput(ABAConfData.Round+1, abaState.GetValues(ABAConfData.Round)[0])

			// if value has only one value, sends FINISH
			if abaState.GetValues(ABAConfData.Round)[0] == s {
				if i.verbose {
					fmt.Println("\tvalue equal to S")
				}
				// check if indeed never sent FINISH message for this vote
				if !abaState.sentFinish(s) {
					finishMsg, err := CreateABAFinish(i.State, i.config, s, ABAConfData.ACRound)
					if err != nil {
						return errors.Wrap(err, "uponABAConf: failed to create ABA Finish message")
					}
					if i.verbose {
						fmt.Println("\tSending ABAFinish")
					}
					i.Broadcast(finishMsg)
					abaState.setSentFinish(s, true)
					abaState.setFinish(i.State.Share.OperatorID, s)
					if i.verbose {
						fmt.Println("\tupdated SentFinish:", abaState.SentFinish)
					}
				}
			}
		}

		// increment round
		if i.verbose {
			fmt.Println("\twill icrement round. Round now:", abaState.Round)
		}
		abaState.IncrementRound()
		if i.verbose {
			fmt.Println("\tnew round:", abaState.Round)
		}

		// start new round sending INIT message with vote
		initMsg, err := CreateABAInit(i.State, i.config, abaState.getVInput(abaState.Round), abaState.Round, ABAConfData.ACRound)
		if err != nil {
			return errors.Wrap(err, "uponABAConf: failed to create ABA Init message")
		}
		if i.verbose {
			fmt.Println("\tSending ABAInit with new Vin:", abaState.Vin[abaState.Round], ", for round:", abaState.Round)
		}
		i.Broadcast(initMsg)
	}

	return nil
}

func isValidABAConf(
	state *State,
	config IConfig,
	signedMsg *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedMsg.Message.MsgType != ABAConfMsgType {
		return errors.New("msg type is not ABAConfMsgType")
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

	ABAConfData, err := signedMsg.Message.GetABAConfData()
	if err != nil {
		return errors.Wrap(err, "could not get ABAConfData data")
	}
	if err := ABAConfData.Validate(); err != nil {
		return errors.Wrap(err, "ABAConfData invalid")
	}

	// vote
	votes := ABAConfData.Votes
	for _, vote := range votes {
		if vote != 0 && vote != 1 {
			return errors.New("vote different than 0 and 1")
		}
	}

	return nil
}

func CreateABAConf(state *State, config IConfig, votes []byte, round Round, acRound ACRound) (*SignedMessage, error) {
	ABAConfData := &ABAConfData{
		Votes:   votes,
		Round:   round,
		ACRound: acRound,
	}
	dataByts, err := ABAConfData.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAConf: could not encode abaconf data")
	}
	msg := &Message{
		MsgType:    ABAConfMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,
		Data:       dataByts,
	}
	sig, err := config.GetSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "CreateABAConf: failed signing abaconf msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   msg,
	}
	return signedMsg, nil
}
