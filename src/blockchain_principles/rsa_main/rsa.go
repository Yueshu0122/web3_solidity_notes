package main

import (
	"fmt"
	"time"
	"web3_solidity_notes/src/blockchain_principles/pow"
	rsa_util "web3_solidity_notes/src/blockchain_principles/rsa"
)

func main() {
	nickname := "yournickname"
	nonce := 0

	// 目标：找到以 5 个 0 开头的哈希值
	startTime := time.Now()
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

	// RSA 签名和验证
	privateKey, publicKey, err := rsa_util.GenerateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	// 对 nickname 和 nonce 的组合进行签名
	message := []byte(fmt.Sprintf("%s%d", nickname, nonce))
	signature, err := rsa_util.SignMessage(privateKey, message)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	// 验证签名
	err = rsa_util.VerifySignature(publicKey, message, signature)
	if err == nil {
		fmt.Println("Signature verified successfully.")
	} else {
		fmt.Println("Failed to verify signature:", err)
	}
}
