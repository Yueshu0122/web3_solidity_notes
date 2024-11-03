package main

import (
	"bytes"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/block", handleBroadcastBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var msg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	mutex.Lock()
	prevBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(prevBlock, msg.BPM)

	if isBlockValid(newBlock, prevBlock) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)

		// 广播新区块给其他节点
		broadcastNewBlock(newBlock)
	}

	mutex.Unlock()

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func broadcastNewBlock(block Block) {
	nodes := []string{"http://node1.example.com", "http://node2.example.com"}
	for _, node := range nodes {
		resp, err := http.Post(node+"/block", "application/json", bytes.NewBuffer(mustMarshal(block)))
		if err != nil {
			log.Printf("Failed to broadcast block to %s: %v", node, err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Node %s rejected the block: %v", node, resp.Status)
		}
	}
}

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}
	return data
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func handleBroadcastBlock(w http.ResponseWriter, r *http.Request) {
	var newBlock Block

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newBlock); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": "Invalid block format"})
		return
	}

	defer r.Body.Close()

	mutex.Lock()
	defer mutex.Unlock()

	if len(Blockchain) == 0 {
		// 如果本地区块链为空，直接添加新区块
		Blockchain = append(Blockchain, newBlock)
		respondWithJSON(w, r, http.StatusOK, map[string]string{"message": "Block added to empty blockchain"})
		return
	}

	prevBlock := Blockchain[len(Blockchain)-1]

	if isBlockValid(newBlock, prevBlock) {
		Blockchain = append(Blockchain, newBlock)
		respondWithJSON(w, r, http.StatusOK, map[string]string{"message": "Block added to blockchain"})
	} else {
		respondWithJSON(w, r, http.StatusBadRequest, map[string]string{"error": "Invalid block"})
	}
}
