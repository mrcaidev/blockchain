package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"
)

// 钱包集数据库。
const walletsDbPath = "wallets.dat"

// 钱包集结构。
type wallets struct {
	Map map[string]*wallet // 钱包地址 - 钱包内容。
}

// 读取钱包集。
func LoadWallets() *wallets {
	// 如果数据库不存在，就返回空钱包集。
	if walletsDbNotExists() {
		return &wallets{make(map[string]*wallet)}
	}

	// 从数据库读取目前的钱包集信息。
	seq, err := ioutil.ReadFile(walletsDbPath)
	if err != nil {
		panic(err)
	}

	return deserializeWallets(seq)
}

// 向钱包集添加钱包。
func (ws *wallets) AddWallet() string {
	wallet := newWallet()
	address := wallet.address()
	ws.Map[address] = wallet
	return address
}

// 获取各钱包的地址。
func (ws *wallets) Addresses() []string {
	var addresses []string
	for address := range ws.Map {
		addresses = append(addresses, address)
	}
	return addresses
}

// 获取指定地址的钱包。
func (ws *wallets) GetWallet(address string) *wallet {
	return ws.Map[address]
}

// 将钱包集存储进数据库。
func (ws *wallets) Persist() {
	err := ioutil.WriteFile(walletsDbPath, ws.serialize(), 0644)
	if err != nil {
		panic(err)
	}
}

// 序列化钱包集。
func (ws *wallets) serialize() []byte {
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
func deserializeWallets(seq []byte) *wallets {
	var wallets wallets

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(seq))
	err := decoder.Decode(&wallets)
	if err != nil {
		panic(err)
	}

	return &wallets
}

// 判断钱包集数据库是否存在。
func walletsDbNotExists() bool {
	_, err := os.Stat(walletsDbPath)
	return os.IsNotExist(err)
}
