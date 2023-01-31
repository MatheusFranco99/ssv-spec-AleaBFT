package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/main/task"
	// spec "github.com/attestantio/go-eth2-client/spec/phase0"
	// "github.com/herumi/bls-eth-go-binary/bls"
)

func getTime() []byte {
	currentTime := time.Now().UnixNano()
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(currentTime))
	return timeBytes
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli launch -n [INT] or propose -n [INT]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "launch":
		n, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Error: argument for -n must be an integer")
			os.Exit(1)
		}
		task.Launch(n)
	case "propose":
		n, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Error: argument for -n must be an integer")
			os.Exit(1)
		}
		task.Propose(n, getTime())
	default:
		fmt.Println("Usage: cli launch -n [INT] or propose -n [INT]")
		os.Exit(1)
	}
}
