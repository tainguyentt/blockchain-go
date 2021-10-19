package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() //load env variables from .env file
	if err != nil {
		log.Fatal(err)
	}

	bcChan = make(chan []Block)

	//init blockchain with a genesis block
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}
	spew.Dump(genesisBlock)

	mutex.Lock()
	Blockchain = append(Blockchain, genesisBlock)
	mutex.Unlock()

	runWebserver()

	// runTcpServer()
}
