package messages

// TODO<olegshmuelov> check if need the invalid tests
// ProposalDataInvalid tests proposal data len == 0
//func ProposalDataInvalid() *tests.MsgSpecTest {
//	d := &qbft.ProposalData{}
//	byts, _ := d.Encode()
//	msg := testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
//		MsgType:    qbft.ProposalMsgType,
//		Height:     qbft.FirstHeight,
//		Round:      10,
//		Identifier: []byte{1, 2, 3, 4},
//		Data:       byts,
//	})
//
//	return &tests.MsgSpecTest{
//		Name: "proposal data invalid",
//		Messages: []*qbft.SignedMessage{
//			msg,
//		},
//		ExpectedError: "ProposalData data is invalid",
//	}
//}
