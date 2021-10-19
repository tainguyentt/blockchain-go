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

	bcServer = make(chan []Block)

	//init blockchain with a genesis block
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	//start webserver
	// log.Fatal((run()))
	runTcpServer()
}
