package tests

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/pkg/errors"
)

var SignedProposal1 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{1, 2, 3, 4}),
})
var SignedProposal2 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{5, 6, 7, 8}),
})
var SignedProposal3 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{1, 3, 5, 7}),
})
var SignedProposal4 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{2, 4, 6, 8}),
})
var SignedProposal5 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{1, 5, 2, 3}),
})
var SignedProposal6 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{1, 6, 2, 3}),
})

var ProposalData1, _ = SignedProposal1.Message.GetProposalData()
var ProposalData2, _ = SignedProposal2.Message.GetProposalData()
var ProposalData3, _ = SignedProposal3.Message.GetProposalData()
var ProposalData4, _ = SignedProposal4.Message.GetProposalData()
var ProposalData5, _ = SignedProposal5.Message.GetProposalData()
var ProposalData6, _ = SignedProposal6.Message.GetProposalData()

var ProposalDataList = []*alea.ProposalData{ProposalData1, ProposalData2}
var ProposalDataList2 = []*alea.ProposalData{ProposalData3, ProposalData4}
var ProposalDataList3 = []*alea.ProposalData{ProposalData5, ProposalData6}

var Entries = [][]*alea.ProposalData{{ProposalData1, ProposalData2}, {ProposalData3, ProposalData4}}
var Priorities = []alea.Priority{alea.FirstPriority, alea.FirstPriority + 1}

var Hash, _ = alea.GetProposalsHash([]*alea.ProposalData{ProposalData1, ProposalData2})
var Hash2, _ = alea.GetProposalsHash([]*alea.ProposalData{ProposalData3, ProposalData4})
var Hash3, _ = alea.GetProposalsHash([]*alea.ProposalData{ProposalData5, ProposalData6})

var AggregatedMsgBytes = func(author types.OperatorID, priority alea.Priority) []byte {

	readyMsgs := make([]*alea.SignedMessage, 0)
	for opID := 1; opID <= 3; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(Hash, priority, author),
		})
		readyMsgs = append(readyMsgs, signedMessage)
	}

	// get ready message and make proof
	aggregatedReadyMessages, err := alea.AggregateMsgs(readyMsgs)
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	aggregatedMsgBytes, err := aggregatedReadyMessages.Encode()
	if err != nil {
		errors.Wrap(err, "could not encode aggregated msg")
	}
	return aggregatedMsgBytes
}

var AggregatedMsgBytes2 = func(author types.OperatorID, priority alea.Priority) []byte {

	readyMsgs := make([]*alea.SignedMessage, 0)
	for opID := 1; opID <= 3; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(Hash2, priority, author),
		})
		readyMsgs = append(readyMsgs, signedMessage)
	}

	// get ready message and make proof
	aggregatedReadyMessages, err := alea.AggregateMsgs(readyMsgs)
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	aggregatedMsgBytes, err := aggregatedReadyMessages.Encode()
	if err != nil {
		errors.Wrap(err, "could not encode aggregated msg")
	}
	return aggregatedMsgBytes
}

var AggregatedMsgBytes3 = func(author types.OperatorID, priority alea.Priority) []byte {

	readyMsgs := make([]*alea.SignedMessage, 0)
	for opID := 1; opID <= 3; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(Hash3, priority, author),
		})
		readyMsgs = append(readyMsgs, signedMessage)
	}

	// get ready message and make proof
	aggregatedReadyMessages, err := alea.AggregateMsgs(readyMsgs)
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	aggregatedMsgBytes, err := aggregatedReadyMessages.Encode()
	if err != nil {
		errors.Wrap(err, "could not encode aggregated msg")
	}
	return aggregatedMsgBytes
}

var AggregatedMsgBytesList = func(author types.OperatorID, priority alea.Priority) [][]byte {
	return [][]byte{AggregatedMsgBytes(author, priority), AggregatedMsgBytes2(author, priority+1)}
}
