package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// models
type Message struct {
	BPM int
}

func runWebserver() {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// router
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// handlers
func handleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	consensus := os.Getenv("CONSENSUS_ALGO")
	var blockchain interface{}
	if consensus == "PoW" {
		blockchain = Blockchain
	} else {
		blockchain = PosBlockchain
	}

	bytes, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJson(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	consensus := os.Getenv("CONSENSUS_ALGO")
	if consensus == "PoS" {
		log.Fatal("PoS consensus is not supported for creating new block via PoS")
	}
	newBlock, err := mineNewPoWBlock(m.BPM)
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, r.Body)
		return
	}
	respondWithJson(w, r, http.StatusCreated, newBlock)
}

func respondWithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(([]byte("HTTP 500: Internal Server Error")))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}
