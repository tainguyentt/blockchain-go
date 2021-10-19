package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
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

// validation
func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// proof of work
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

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
		newBlockHash := calculateHash(newBlock)
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

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write(([]byte(record)))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
