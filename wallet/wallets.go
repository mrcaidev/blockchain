package wallet

import (
	"blockchain/utils"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

// 钱包集数据库。
const walletsPath = "wallets.dat"

// 钱包集结构。
type Wallets struct {
	Wallets map[string]*Wallet
}

// 创建钱包集。
func NewWallets() *Wallets {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	wallets.LoadFromFile()
	return &wallets
}

// 读取钱包集。
func (wallets *Wallets) LoadFromFile() {
	// 如果没有钱包数据库，就给出提示。
	if !utils.HasFile(walletsPath) {
		fmt.Println("wallets.dat not found.")
		return
	}

	// 读取数据库。
	seq, err := ioutil.ReadFile(walletsPath)
	if err != nil {
		panic(err)
	}
	wallets.Wallets = DeserializeWallets(seq).Wallets
}

// 向钱包集添加钱包。
func (wallets *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := wallet.Address()
	wallets.Wallets[address] = wallet
	return address
}

// 获取各钱包的地址。
func (wallets *Wallets) Addresses() []string {
	var addresses []string
	for addr := range wallets.Wallets {
		addresses = append(addresses, addr)
	}
	return addresses
}

// 获取指定钱包的地址。
func (wallets *Wallets) GetWallet(address string) *Wallet {
	return wallets.Wallets[address]
}

// 将钱包集存储进数据库。
func (wallets *Wallets) Store() {
	err := ioutil.WriteFile(walletsPath, wallets.Serialize(), 0644)
	if err != nil {
		panic(err)
	}
}

// 序列化钱包集。
func (wallets *Wallets) Serialize() []byte {
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
func DeserializeWallets(seq []byte) *Wallets {
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(seq))
	err := decoder.Decode(&wallets)
	if err != nil {
		panic(err)
	}
	return &wallets
}
