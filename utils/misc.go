package utils

import (
	"bytes"
	"crypto/sha256"
)

// 计算校验和。
func Checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:ChecksumLen]
}

// 检验地址正确性。
func ValidateAddress(address string) bool {
	pubkeyHash := Base58Decode([]byte(address))
	actualChecksum := pubkeyHash[len(pubkeyHash)-ChecksumLen:]
	version := pubkeyHash[0]
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash)-ChecksumLen]
	supposedChecksum := Checksum(append([]byte{version}, pubkeyHash...))
	return bytes.Equal(actualChecksum, supposedChecksum)
}
