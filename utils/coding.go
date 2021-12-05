package utils

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

var alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58编码。
func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}

	result = reverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{alphabet[0]}, result...)
		} else {
			break
		}
	}
	return result
}

// Base58解码。
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}

// 翻转byte数据。
func reverseBytes(input []byte) []byte {
	var result []byte
	for i := len(input) - 1; i >= 0; i-- {
		result = append(result, input[i])
	}
	return result
}

// 计算公钥的哈希值。
func HashPubKey(pubkey []byte) []byte {
	pubSHA := sha256.Sum256(pubkey)
	hasher := ripemd160.New()
	_, err := hasher.Write(pubSHA[:])
	if err != nil {
		log.Panic(err)
	}
	return hasher.Sum(nil)
}
