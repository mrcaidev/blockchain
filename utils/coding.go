package utils

import (
	"bytes"
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// Base58 字母表。
var alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58 基数。
const base = 58

// 校验和字节长度。
const checksumLen = 4

// Base58 编码。
func Base58Encode(input []byte) []byte {
	output := []byte{}
	bigInput := BytesToBigInt(input)
	bigBase := big.NewInt(int64(base))
	bigZero := big.NewInt(0)
	mod := big.NewInt(0)

	for bigInput.Cmp(bigZero) != 0 {
		bigInput.DivMod(bigInput, bigBase, mod)
		output = append(output, alphabet[mod.Int64()])
	}

	output = reverseBytes(output)
	for b := range input {
		if b == 0x00 {
			output = append([]byte{alphabet[0]}, output...)
		} else {
			break
		}
	}
	return output
}

// Base58 解码。
func Base58Decode(input []byte) []byte {
	output := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(alphabet, b)
		output.Mul(output, big.NewInt(base))
		output.Add(output, big.NewInt(int64(charIndex)))
	}

	decoded := output.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}

// 首尾翻转 []byte 类型的数据。
func reverseBytes(input []byte) []byte {
	output := []byte{}
	for pos := len(input) - 1; pos >= 0; pos-- {
		output = append(output, input[pos])
	}
	return output
}

// 计算公钥的哈希值。（SHA256 + RIPEMD）
func GetPubkeyHash(pubkey []byte) []byte {
	first := sha256.Sum256(pubkey)
	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(first[:])
	if err != nil {
		panic(err)
	}
	return ripemdHasher.Sum(nil)
}

// 计算校验和。（两次 SHA256）
func GetChecksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:checksumLen]
}
