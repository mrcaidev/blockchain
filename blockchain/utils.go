package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
)

// 将整型数转换为字节切片。
func IntToBytes(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// 将字节切片转换为整型数。
func BytesToBigInt(bytes []byte) *big.Int {
	var result big.Int
	result.SetBytes(bytes)
	return &result
}
