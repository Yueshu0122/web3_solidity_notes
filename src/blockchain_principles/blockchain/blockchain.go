package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var Blockchain []Block
var mutex = &sync.Mutex{}

func syncBlockchain(nodes []string) {
	longestChain := Blockchain
	longestLength := len(Blockchain)

	for _, node := range nodes {
		resp, err := http.Get(node + "/")
		if err != nil {
			fmt.Printf("Failed to fetch blockchain from %s: %v\n", node, err)
			continue
		}
		defer resp.Body.Close()

		var remoteBlockchain []Block
		if err := json.NewDecoder(resp.Body).Decode(&remoteBlockchain); err != nil {
			fmt.Printf("Failed to decode blockchain from %s: %v\n", node, err)
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
		fmt.Println("Blockchain updated to the longest chain")
	}
}

func isChainValid(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		if !isBlockValid(chain[i], chain[i-1]) {
			return false
		}
	}
	return true
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if !isValidProofOfWork(newBlock) {
		return false
	}

	return true
}

func isValidProofOfWork(block Block) bool {
	hash := calculateHashWithNonce(block, block.Nonce)
	return hash[:difficulty] == strings.Repeat("0", difficulty)
}

func slicesEqual(a, b []Block) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
