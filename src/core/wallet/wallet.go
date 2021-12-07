package wallet

import (
	"blockchain/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

// 当前版本号。
const version = byte(0x00)

// 钱包结构。
type wallet struct {
	Privkey ecdsa.PrivateKey // 私钥。
	Pubkey  []byte           // 公钥。
}

// 创建钱包。
func newWallet() *wallet {
	privkey, pubkey := newKeyPair()
	return &wallet{privkey, pubkey}
}

// 获取钱包地址。
// 算法：地址 = (版本号 + 公钥哈希 + 校验和) 的 Base58 编码。
func (w *wallet) address() string {
	pubkeyHash := utils.GetPubkeyHash(w.Pubkey)
	payload := append([]byte{version}, pubkeyHash...)
	checksum := utils.GetChecksum(payload)
	payload = append(payload, checksum...)
	return string(utils.Base58Encode(payload))
}

// 创建新公钥-私钥对。
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 椭圆加密产生私钥。
	curve := elliptic.P256()
	privkey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	// 由私钥衍生出公钥。
	pubkey := append(privkey.PublicKey.X.Bytes(), privkey.Y.Bytes()...)

	return *privkey, pubkey
}
