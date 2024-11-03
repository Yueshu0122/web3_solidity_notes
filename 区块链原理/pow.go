package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// 计算并返回给定数据的 SHA-256 哈希值
func calculateSHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// 检查哈希值是否以指定数量的 0 开头
func hasLeadingZeros(hash string, count int) bool {
	return len(hash) >= count && hash[:count] == strings.Repeat("0", count)
}

func main() {
	nickname := "yournickname"
	nonce := 0

	// 目标：找到以 4 个 0 开头的哈希值
	startTime := time.Now()
	for {
		data := fmt.Sprintf("%s%d", nickname, nonce)
		hash := calculateSHA256([]byte(data))

		if hasLeadingZeros(hash, 4) {
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
		hash := calculateSHA256([]byte(data))

		if hasLeadingZeros(hash, 5) {
			elapsedTime := time.Since(startTime)
			fmt.Printf("Found hash with 5 leading zeros: %s (Nonce: %d, Time: %v)\n", hash, nonce, elapsedTime)
			break
		}
		nonce++
	}
}
