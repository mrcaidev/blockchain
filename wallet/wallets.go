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
	Wallets map[string]*wallet // 钱包地址 - 钱包内容。
}

// 创建钱包集。
func LoadWallets() *Wallets {
	// 如果数据库不存在，就返回空钱包集。
	if !utils.HasFile(walletsPath) {
		return &Wallets{make(map[string]*wallet)}
	}

	// 从数据库读取目前的钱包集信息。
	seq, err := ioutil.ReadFile(walletsPath)
	if err != nil {
		panic(err)
	}

	return deserializeWallets(seq)
}

// 向钱包集添加钱包。
func (ws *Wallets) AddWallet() string {
	wallet := newWallet()
	address := wallet.address()
	ws.Wallets[address] = wallet
	return address
}

// 获取各钱包的地址。
func (ws *Wallets) Addresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// 获取指定地址的钱包。
func (ws *Wallets) GetWallet(address string) *wallet {
	return ws.Wallets[address]
}

// 将钱包集存储进数据库。
func (ws *Wallets) Persist() {
	err := ioutil.WriteFile(walletsPath, ws.serialize(), 0644)
	if err != nil {
		panic(err)
	}
}

// 序列化钱包集。
func (ws *Wallets) serialize() []byte {
	var seq bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&seq)
	err := encoder.Encode(ws)
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
