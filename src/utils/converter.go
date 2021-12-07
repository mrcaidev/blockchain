package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

// int64 -> []byte
func Int64ToBytes(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

// []byte -> big.Int
func BytesToBigInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(data)
}
