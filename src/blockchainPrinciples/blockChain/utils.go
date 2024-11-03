package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func generateBlock(oldBlock Block, BPM int, transactions []Transaction) Block {
	t := time.Now()

	newBlock := Block{
		Index:        oldBlock.Index + 1,
		Timestamp:    t.String(),
		BPM:          BPM,
		PrevHash:     oldBlock.Hash,
		Transactions: transactions,
	}
	nonce, hash := proofOfWork(newBlock)
	newBlock.Nonce = nonce
	newBlock.Hash = hash

	return newBlock
}

func syncBlockchain(nodes []string) {
	longestChain := Blockchain
	longestLength := len(Blockchain)

	for _, node := range nodes {
		resp, err := http.Get(node + "/")
		if err != nil {
			log.Printf("Failed to fetch blockChain from %s: %v\n", node, err)
			continue
		}
		defer resp.Body.Close()

		var remoteBlockchain []Block
		if err := json.NewDecoder(resp.Body).Decode(&remoteBlockchain); err != nil {
			log.Printf("Failed to decode blockChain from %s: %v\n", node, err)
			continue
		}

		if len(remoteBlockchain) > longestLength && isChainValid(remoteBlockchain) {
			longestChain = remoteBlockchain
			longestLength = len(remoteBlockchain)
		}
	}

	if !slicesEqual(longestChain, Blockchain) {
		mutex.Lock()
		Blockchain = longestChain
		mutex.Unlock()
		log.Println("Blockchain updated to the longest chain")
	}
}
