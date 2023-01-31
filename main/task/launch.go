package main

import (

	// spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"

	// "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/main/propose"
)

func getTime() []byte {
	currentTime := time.Now().Unix()
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(currentTime))
	return timeBytes
}

func main() {

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error reading input n:", err)
		return
	}
	inst := testingutils.BaseInstanceAleaN(types.OperatorID(n))
	inst.Start([]byte{1, 2, 3, 4}, alea.FirstHeight)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		input := strings.Split(text, " ")
		switch input[0] {
		case "propose":
			if len(input) == 2 {
				num, err := strconv.Atoi(input[1])
				if err != nil {
					fmt.Println("Invalid argument.")
					continue
				}
				for idx := 1; idx <= num; idx++ {
					propose.Propose(n, getTime())
				}
				// fmt.Printf("Proposing %d...\n", num)
			} else {
				fmt.Println("Invalid command.")
			}
		case "startABA":
			go inst.StartAgreementComponent()
			// fmt.Println("Starting ABA...")
		case "stop":
			inst.State.StopAgreement = true
		default:
			fmt.Println("Invalid command.")
		}
	}

}
