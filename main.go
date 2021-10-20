package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() //load env variables from .env file
	if err != nil {
		log.Fatal(err)
	}

	bcChan = make(chan []Block)

	var connHandler func(net.Conn)
	consensusAlgo := os.Getenv("CONSENSUS_ALGO")
	switch consensusAlgo {
	case "PoW":
		connHandler = handlePowConn
		initPowGenesisBlock()
		break
	case "PoS":
		connHandler = handlePosConn
		initPosGenesisBlock()

		//start a goroutine to add incoming blocks
		go func() {
			for candidate := range candidateBlocks {
				mutex.Lock()
				tempPosBlocks = append(tempPosBlocks, candidate)
				mutex.Unlock()
			}
		}()

		//start a goroutine to pick a winner and mint a new block
		go func() {
			for {
				pickWinner()
			}
		}()
		break
	default:
		log.Fatal("undefined consensus algorithm")
	}

	//start webserver
	go runWebserver()

	// start tcp socket
	runTcpServer(connHandler)
}
