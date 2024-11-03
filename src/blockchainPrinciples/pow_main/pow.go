package main

import (
	"fmt"
	"time"
	"web3_solidity_notes/src/blockchainPrinciples/pow"
)

func main() {
	nickname := "yournickname"
	nonce := 0

	// 目标：找到以 4 个 0 开头的哈希值
	startTime := time.Now()
	for {
		data := fmt.Sprintf("%s%d", nickname, nonce)
		hash := pow.CalculateSHA256([]byte(data))

		if pow.HasLeadingZeros(hash, 4) {
			elapsedTime := time.Since(startTime)
			fmt.Printf("Found hash with 4 leading zeros: %s (Nonce: %d, Time: %v)\n", hash, nonce, elapsedTime)
			break
		}
		nonce++
	}

	// 目标：找到以 5 个 0 开头的哈希值
	startTime = time.Now()
	for {
		data := fmt.Sprintf("%s%d", nickname, nonce)
		hash := pow.CalculateSHA256([]byte(data))

		if pow.HasLeadingZeros(hash, 5) {
			elapsedTime := time.Since(startTime)
			fmt.Printf("Found hash with 5 leading zeros: %s (Nonce: %d, Time: %v)\n", hash, nonce, elapsedTime)
			break
		}
		nonce++
	}
}
