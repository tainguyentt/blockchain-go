package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type PosBlock struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Validator string
}

var PosBlockchain []PosBlock
var tempPosBlocks []PosBlock

// channel to handle incoming blocks to be validated
var candidateBlocks = make(chan PosBlock)

// channel to broadcast winning validator to all nodes
var announcements = make(chan string)

// control read/write to prevent data race
var mu sync.Mutex

//keep track of validators and their balances
var validators = make(map[string]int)

// utility
func calculatePosBlockHash(block PosBlock) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	return hash(record)
}

// mint a new PoS block
func generatePosBlock(oldBlock PosBlock, BPM int, address string) (PosBlock, error) {
	var newBlock PosBlock
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculatePosBlockHash(newBlock)
	newBlock.Validator = address
	return newBlock, nil
}

func initPosGenesisBlock() {
	t := time.Now()
	genesisBlock := PosBlock{}
	genesisBlock = PosBlock{0, t.String(), 0, calculatePosBlockHash(genesisBlock), "", ""}
	spew.Dump(genesisBlock)
	PosBlockchain = append(PosBlockchain, genesisBlock)
}

// validate block
func isPosBlockValid(newBlock PosBlock, oldBlock PosBlock) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculatePosBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// add a new node to validators, balance is its total stake
func addValidator(address string, balance int) {
	validators[address] = balance
	fmt.Println(validators)
}

// pick a winner
func pickWinner() {
	fmt.Println("Picking up a winner")
	time.Sleep(30 * time.Second)

	mu.Lock()
	tempBlocks := tempPosBlocks
	mu.Unlock()

	lotteryPool := []string{}
	if len(tempBlocks) > 0 {
	OUTER:
		for _, block := range tempBlocks {
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			//a validator with more stakes will have more lots in the pool
			mu.Lock()
			validatorSet := validators
			mu.Unlock()
			k, ok := validatorSet[block.Validator]
			if ok {
				for i := 0; i < k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// randomly pick a winner from the pool
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		winner := lotteryPool[r.Intn(len(lotteryPool))]

		// add block of winner to blockchain and let all other nodes know
		for _, block := range tempBlocks {
			if block.Validator == winner {
				mu.Lock()
				PosBlockchain = append(PosBlockchain, block)
				mu.Unlock()

				for range validators {
					announceMsg := "\nWinning validator: " + winner + "\n"
					announcements <- announceMsg
				}
				break
			}
		}

		// reset temp blocks
		mu.Lock()
		tempPosBlocks = []PosBlock{}
		mu.Unlock()
	}
}

func handlePosConn(conn net.Conn) {
	defer conn.Close()

	//goroutine to broadcast new blocks to other peers
	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()

	var address string

	// validators to input their addresses and staking tokens
	io.WriteString(conn, "Enter a token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = hash(t.String())
		addValidator(address, balance)
		break
	}

	//allow connected node to continuously input their BPM data
	io.WriteString(conn, "\nEnter a new BPM:")
	scanBPM := bufio.NewScanner(conn)
	go func() {
		for {
			for scanBPM.Scan() {
				bpm, err := strconv.Atoi(scanBPM.Text())
				// terminate connection if invalid data
				if err != nil {
					log.Printf("%v is not a number: %v", scanBPM.Text(), err)
					delete(validators, address)
					conn.Close()
				}

				mutex.Lock()
				oldLastIndex := PosBlockchain[len(PosBlockchain)-1]
				mutex.Unlock()

				newBlock, err := generatePosBlock(oldLastIndex, bpm, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isPosBlockValid(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new BPM:")
			}
		}
	}()

	// simulate receiving broadcast message
	for {
		time.Sleep(time.Minute)

		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()

		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(conn, string(output)+"\n")
	}
}
