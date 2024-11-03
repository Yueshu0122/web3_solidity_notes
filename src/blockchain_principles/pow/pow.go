package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// 计算并返回给定数据的 SHA-256 哈希值
func CalculateSHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// 检查哈希值是否以指定数量的 0 开头
func HasLeadingZeros(hash string, count int) bool {
	return len(hash) >= count && hash[:count] == strings.Repeat("0", count)
}
