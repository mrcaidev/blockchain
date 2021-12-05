package wallet

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

// 钱包结构。
type Wallet struct {
	PrivateKey ecdsa.PrivateKey // 私钥。
	PublicKey  []byte           // 公钥。
}

// 创建钱包。
func NewWallet() *Wallet {
	privkey, pubkey := newKeyPair()
	return &Wallet{privkey, pubkey}
}

// 获取钱包地址。
func (wallet *Wallet) Address() string {
	pubkeyHash := utils.GetPubkeyHash(wallet.PublicKey)
	payloadWithVersion := append([]byte{utils.Version}, pubkeyHash...)
	checksum := utils.GetChecksum(payloadWithVersion)
	finalPayload := append(payloadWithVersion, checksum...)
	return string(utils.Base58Encode(finalPayload))
}

// 创建新公钥-私钥对。
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privkey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	pubkey := append(privkey.PublicKey.X.Bytes(), privkey.Y.Bytes()...)
	return *privkey, pubkey
}
