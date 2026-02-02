package security

import (
	"crypto/sha1"
	"fmt"
)

// Sha1 计算 SHA1 哈希值
func Sha1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ValidateCheckSum 验证 checksum 是否有效
// checksum = SHA1(secret + nonce + timestamp)
func ValidateCheckSum(checksum, timestamp, nonce, secret string) bool {
	calculatedSum := Sha1(secret + nonce + timestamp)
	return calculatedSum == checksum
}
