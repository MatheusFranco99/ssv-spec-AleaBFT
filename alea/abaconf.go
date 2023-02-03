package alea

import (
	"fmt"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/pkg/errors"
)

func (i *Instance) uponABAConf(signedABAConf *SignedMessage) error {
	if i.verbose {
		fmt.Println(Teal("#######################################################"))
		fmt.Println(Teal("uponABAConf"))
	}
	// get data
	ABAConfData, err := signedABAConf.Message.GetABAConfData()
	if err != nil {
		return errors.Wrap(err, "uponABAConf:could not get ABAConfData from signedABAConf")
	}

	if i.verbose {
		fmt.Println(Teal("\tACRound:", ABAConfData.ACRound, "Round:", ABAConfData.Round, "Votes:", ABAConfData.Votes))
		fmt.Println(Teal("\tOwn ACState.ACRound:", i.State.ACState.ACRound))
	}

	// old message -> ignore
	if ABAConfData.ACRound < i.State.ACState.ACRound {
		if i.verbose {
			fmt.Println(Teal("\told message. Returning..."))
		}
		return nil
	}
	if ABAConfData.ACRound == i.State.ACState.ACRound && ABAConfData.Round < i.State.ACState.GetCurrentABAState().Round {
		if i.verbose {
			fmt.Println(Teal("\told message. Returning..."))
		}
		return nil
	}

	// if future round -> intialize future state
	if ABAConfData.ACRound > i.State.ACState.ACRound {
		i.State.ACState.InitializeRound(ABAConfData.ACRound)
	}
	if ABAConfData.Round > i.State.ACState.GetABAState(ABAConfData.ACRound).Round {
		i.State.ACState.GetABAState(ABAConfData.ACRound).InitializeRound(ABAConfData.Round)
	}

	abaState := i.State.ACState.GetABAState(ABAConfData.ACRound)

	// add the message to the containers
	abaState.ABAConfContainer.AddMsg(signedABAConf)

	// sender
	senderID := signedABAConf.GetSigners()[0]

	alreadyReceived := abaState.hasConf(ABAConfData.Round, senderID)
	if i.verbose {
		fmt.Println(Teal("\tsenderID:", senderID, ", already received before:", alreadyReceived))
	}
	// if never received this msg, update
	if !alreadyReceived {

		// determine if votes list is contained in local round values list
		isContained := abaState.isContainedInValues(ABAConfData.Round, ABAConfData.Votes)

		if i.verbose {
			fmt.Println(Teal("\tis value contained in own values? ", isContained))
		}

		// list is contained -> update CONF counter
		if isContained {
			abaState.setConf(ABAConfData.Round, senderID)
			if i.verbose {
				fmt.Println(Teal("\tupdated confcounter:", abaState.countConf(ABAConfData.Round)))
			}
		}
	}

	// reached strong support -> try to decide value
	if abaState.countConf(ABAConfData.Round) >= i.State.Share.Quorum {
		if i.verbose {
			fmt.Println(Teal("\treached quorum of conf"))
		}

		// get common coin
		s := abaState.Coin(abaState.Round)
		if i.verbose {
			fmt.Println(Teal("\tcoin:", s))
		}

		// if values = {0,1}, choose randomly (i.e. coin) value for next round
		if len(abaState.Values[ABAConfData.Round]) == 2 {

			abaState.setVInput(ABAConfData.Round+1, s)
			if i.verbose {
				fmt.Println(Teal("\tlength of values is 2", abaState.Values[ABAConfData.Round], "-> storing coin to next Vin"))
			}
		} else {
			if i.verbose {
				fmt.Println(Teal("\tlength of values is 1:", abaState.Values[ABAConfData.Round]))
			}
			abaState.setVInput(ABAConfData.Round+1, abaState.GetValues(ABAConfData.Round)[0])

			// if value has only one value, sends FINISH
			if abaState.GetValues(ABAConfData.Round)[0] == s {
				if i.verbose {
					fmt.Println(Teal("\tvalue equal to S"))
				}
				// check if indeed never sent FINISH message for this vote
				if !abaState.sentFinish(s) {
					finishMsg, err := CreateABAFinish(i.State, i.config, s, ABAConfData.ACRound)
					if err != nil {
						return errors.Wrap(err, "uponABAConf: failed to create ABA Finish message")
					}
					if i.verbose {
						fmt.Println(Teal("\tSending ABAFinish"))
					}
					i.Broadcast(finishMsg)
					abaState.setSentFinish(s, true)
					abaState.setFinish(i.State.Share.OperatorID, s)
					if i.verbose {
						fmt.Println(Teal("\tupdated SentFinish:", abaState.SentFinish))
					}
				}
			}
		}

		// increment round
		if i.verbose {
			fmt.Println(Teal("\twill increment round. Current round:", abaState.Round))
		}
		abaState.IncrementRound()
		if i.verbose {
			fmt.Println(Teal("\tNew round:", abaState.Round))
		}

		// start new round sending INIT message with vote
		initMsg, err := CreateABAInit(i.State, i.config, abaState.getVInput(abaState.Round), abaState.Round, ABAConfData.ACRound)
		if err != nil {
			fmt.Println(Teal(err, "uponABAConf: failed to create ABA Init message"))
			return errors.Wrap(err, "uponABAConf: failed to create ABA Init message")
		}
		if i.verbose {
			fmt.Println(Teal("\tSending ABAInit with new Vin:", abaState.Vin[abaState.Round], ", for round:", abaState.Round, ", for ACRound:", ABAConfData.ACRound))
		}
		i.Broadcast(initMsg)

		abaState.setSentInit(abaState.Round, abaState.getVInput(abaState.Round), true)
		abaState.setInit(abaState.Round, i.State.Share.OperatorID, abaState.getVInput(abaState.Round))

		if i.verbose {
			fmt.Println(Teal("\tupdating own initCounter and setInit:", abaState.InitCounter, ", ", abaState.SentInit))
		}

		if i.verbose {
			fmt.Println(Teal("finishABAConf"))
			fmt.Println(Teal("#######################################################"))
		}
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
