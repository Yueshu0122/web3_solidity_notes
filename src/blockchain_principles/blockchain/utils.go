package main

import (
	"time"
)

func generateBlock(oldBlock Block, BPM int) Block {
	t := time.Now()

	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: t.String(),
		BPM:       BPM,
		PrevHash:  oldBlock.Hash,
	}
	nonce, hash := proofOfWork(newBlock)
	newBlock.Nonce = nonce
	newBlock.Hash = hash

	return newBlock
}
