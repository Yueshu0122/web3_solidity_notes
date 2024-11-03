package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index        int
	Timestamp    string
	BPM          int
	Hash         string
	PrevHash     string
	Nonce        int
	Transactions []Transaction
}

type Transaction struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Amount    int    `json:"amount"`
	Timestamp string `json:"timestamp"`
}

var nodeList = []string{
	"http://localhost:8001",
	"http://localhost:8002",
	"http://localhost:8003",
	"http://localhost:8004",
	"http://localhost:8005",
}

const difficulty = 4 // 难度值，表示哈希值前缀需要多少个0

func NewBlock(index int, timestamp time.Time, bpm int, prevHash string, transactions []Transaction) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    timestamp.String(),
		BPM:          bpm,
		PrevHash:     prevHash,
		Transactions: transactions,
	}
	nonce, hash := proofOfWork(*block)
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func proofOfWork(block Block) (int, string) {
	nonce := 0
	for {
		nonce++
		hash := calculateHashWithNonce(block, nonce)
		if hash[:difficulty] == strings.Repeat("0", difficulty) {
			return nonce, hash
		}
	}
}

func calculateHashWithNonce(block Block, nonce int) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + strconv.Itoa(nonce)
	for _, tx := range block.Transactions {
		record += tx.Sender + tx.Receiver + strconv.Itoa(tx.Amount) + tx.Timestamp
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
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
		if !blocksEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func blocksEqual(a, b Block) bool {
	if a.Index != b.Index {
		return false
	}
	if a.Timestamp != b.Timestamp {
		return false
	}
	if a.BPM != b.BPM {
		return false
	}
	if a.Hash != b.Hash {
		return false
	}
	if a.PrevHash != b.PrevHash {
		return false
	}
	if a.Nonce != b.Nonce {
		return false
	}
	if !transactionsEqual(a.Transactions, b.Transactions) {
		return false
	}
	return true
}

func transactionsEqual(a, b []Transaction) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Sender != b[i].Sender ||
			a[i].Receiver != b[i].Receiver ||
			a[i].Amount != b[i].Amount ||
			a[i].Timestamp != b[i].Timestamp {
			return false
		}
	}
	return true
}
