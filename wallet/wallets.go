package wallet

import (
	"blockchain/utils"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
)

// 钱包集数据库。
const walletsPath = "wallets.dat"

// 钱包集结构。
type Wallets struct {
	Wallets map[string]*Wallet // 钱包地址 - 钱包内容。
}

// 创建钱包集。
func LoadWallets() *Wallets {
	// 如果数据库不存在，就返回空钱包集。
	if !utils.HasFile(walletsPath) {
		return &Wallets{make(map[string]*Wallet)}
	}

	// 从数据库读取目前的钱包集信息。
	seq, err := ioutil.ReadFile(walletsPath)
	if err != nil {
		panic(err)
	}

	return deserializeWallets(seq)
}

// 向钱包集添加钱包。
func (wallets *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := wallet.address()
	wallets.Wallets[address] = wallet
	return address
}

// 获取各钱包的地址。
func (wallets *Wallets) Addresses() []string {
	var addresses []string
	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// 获取指定地址的钱包。
func (wallets *Wallets) GetWallet(address string) *Wallet {
	return wallets.Wallets[address]
}

// 将钱包集存储进数据库。
func (wallets *Wallets) Persist() {
	err := ioutil.WriteFile(walletsPath, wallets.serialize(), 0644)
	if err != nil {
		panic(err)
	}
}

// 序列化钱包集。
func (wallets *Wallets) serialize() []byte {
	var seq bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&seq)
	err := encoder.Encode(wallets)
	if err != nil {
		panic(err)
	}

	return seq.Bytes()
}

// 反序列化钱包集。
func deserializeWallets(seq []byte) *Wallets {
	var wallets Wallets

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(seq))
	err := decoder.Decode(&wallets)
	if err != nil {
		panic(err)
	}

	return &wallets
}
