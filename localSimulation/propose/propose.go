package propose

import (

	// spec "github.com/attestantio/go-eth2-client/spec/phase0"

	"fmt"
	"net"
	"strconv"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/pkg/errors"

	// "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
)

func SendTCP(msg *alea.SignedMessage, operatorID types.OperatorID) error {
	byts, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode message")
	}

	// return i.config.GetNetwork().Broadcast(msgToBroadcast)
	port := strconv.Itoa(8000 + int(operatorID))
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println(err)
	} else {
		_, write_err := conn.Write(byts)
		if write_err != nil {
			fmt.Println("failed:", write_err)
		}
		conn.(*net.TCPConn).CloseWrite()
	}
	return nil
}

func Propose(n int, data []byte) {
	// n, err := strconv.Atoi(os.Args[1])
	// if err != nil {
	// 	fmt.Println("Error reading input n:", err)
	// 	return
	// }

	signedMessage := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(n)], types.OperatorID(n), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(data),
	})

	err := SendTCP(signedMessage, types.OperatorID(n))
	if err != nil {
		fmt.Println(err)
	}
}
