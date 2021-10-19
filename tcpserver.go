package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func runTcpServer() {
	tcpPort := os.Getenv("ADDR")
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP  Server Listening on port :", tcpPort)

	defer server.Close()

	// infinite loop to accept new connections
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// new go routine
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}

			lastBlock := Blockchain[len(Blockchain)-1]
			newBlock, err := generateBlock(lastBlock, bpm)
			if err != nil {
				log.Println(err)
				continue
			}

			if isBlockValid(newBlock, lastBlock) {
				newBlockChain := append(Blockchain, newBlock)
				replaceChain(newBlockChain)
			}

			bcChan <- Blockchain
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	//simulating receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range bcChan {
		spew.Dump(Blockchain)
	}
}
