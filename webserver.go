package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	bytes, err := json.MarshalIndent(Blockchain, "", " ") //todo
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	io.WriteString(w, string(bytes)) //todo
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJson(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	lastBlock := Blockchain[len(Blockchain)-1]
	newBlock, err := generateBlock(lastBlock, m.BPM)
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, r.Body)
		return
	}

	if isBlockValid(newBlock, lastBlock) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
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
