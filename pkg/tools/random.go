package tools

import (
	"crypto/rand"
	"math/big"
)

const (
	// 字符集定义
	alphaLower   = "abcdefghijklmnopqrstuvwxyz"
	alphaUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers      = "0123456789"
	alphaNumeric = alphaLower + alphaUpper + numbers
)

// RandString 生成指定长度的随机字符串（包含大小写字母和数字）
func RandString(length int) string {
	return randStringFromCharset(length, alphaNumeric)
}

// RandStringLower 生成指定长度的小写字母随机字符串
func RandStringLower(length int) string {
	return randStringFromCharset(length, alphaLower)
}

// RandStringUpper 生成指定长度的大写字母随机字符串
func RandStringUpper(length int) string {
	return randStringFromCharset(length, alphaUpper)
}

// RandNumber 生成指定长度的数字随机字符串
func RandNumber(length int) string {
	return randStringFromCharset(length, numbers)
}

// RandStringCustom 使用自定义字符集生成随机字符串
func RandStringCustom(length int, charset string) string {
	return randStringFromCharset(length, charset)
}

// randStringFromCharset 从指定字符集生成随机字符串
func randStringFromCharset(length int, charset string) string {
	if length <= 0 {
		return ""
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomInt, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// 如果随机数生成失败，使用字符集的第一个字符
			result[i] = charset[0]
			continue
		}
		result[i] = charset[randomInt.Int64()]
	}

	return string(result)
}
