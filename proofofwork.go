package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const difficulty = 1

// models
type Block struct {
	Index      int    //position of block in the blockchain
	Timestamp  string //created time
	BPM        int    //beats per minute/pulse rate
	Hash       string //SHA256 hash of this block
	PrevHash   string //SHA256 of previous block
	Difficulty int    //the number of leading 0s of the next block hash
	Nonce      string //answer to the math problem needed to solve
}

var Blockchain []Block

var bcChan chan []Block

var mutex sync.Mutex

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// proof of work
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// mine a new block using PoW
func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		newBlockHash := calculatePowBlockHash(newBlock)
		if !isHashValid(newBlockHash, newBlock.Difficulty) {
			fmt.Println(newBlockHash, " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(newBlockHash, " work done!")
			newBlock.Hash = newBlockHash
			break
		}
	}

	return newBlock, nil
}

func initPowGenesisBlock() {
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculatePowBlockHash(genesisBlock), "", difficulty, ""}
	spew.Dump(genesisBlock)

	mutex.Lock()
	Blockchain = append(Blockchain, genesisBlock)
	mutex.Unlock()
}

func calculatePowBlockHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	return hash(record)
}

// validation
func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculatePowBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// handle incoming connection
func handlePowConn(conn net.Conn) {
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
