package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
	Nonce     int
}

type Message struct {
	BPM int
}

const difficulty = 4 // 难度值，表示哈希值前缀需要多少个0

func NewBlock(index int, timestamp time.Time, bpm int, prevHash string) *Block {
	block := &Block{
		Index:     index,
		Timestamp: timestamp.String(),
		BPM:       bpm,
		PrevHash:  prevHash,
	}
	nonce, hash := proofOfWork(*block)
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
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
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
