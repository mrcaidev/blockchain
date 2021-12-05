package utils

import (
	"bytes"
	"os"
)

// 判断是否是有效地址。
func IsValidAddress(address string) bool {
	pubkeyHash := Base58Decode([]byte(address))
	actualChecksum := pubkeyHash[len(pubkeyHash)-checksumLen:]
	supposedChecksum := GetChecksum(pubkeyHash[:len(pubkeyHash)-checksumLen])
	return bytes.Equal(actualChecksum, supposedChecksum)
}

// 检测文件是否存在。
func HasFile(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
