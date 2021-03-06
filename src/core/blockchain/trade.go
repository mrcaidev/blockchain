package blockchain

import (
	"blockchain/core/transaction"
	"blockchain/core/wallet"
	"blockchain/utils"
	"encoding/hex"
	"fmt"
)

// 创建一笔 coinbase 交易。
func NewCoinbaseTx(to string, data string) *transaction.Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// 创建交易的输入和输出。
	txi := transaction.NewTxi([]byte{}, -1, nil, []byte(data))
	txo := transaction.NewTxo(subsidy, to)
	tx := transaction.Transaction{
		ID:      nil,
		Inputs:  []*transaction.TxInput{txi},
		Outputs: []*transaction.TxOutput{txo},
	}
	tx.ID = tx.Hash()
	return &tx
}

// 创建一笔 UTXO 交易。
func (c *Chain) NewUtxoTx(from string, to string, amount int) *transaction.Transaction {
	var (
		newInputs  []*transaction.TxInput
		newOutputs []*transaction.TxOutput
	)

	// 获取发起方的钱包。
	wallets := wallet.LoadWallets()
	wallet := wallets.GetWallet(from)
	pubkeyHash := utils.GetPubkeyHash(wallet.Pubkey)

	// 从钱包里找出足够多的钱。
	deposit, UTXOToPay := c.FindUtxosToPay(pubkeyHash, amount)

	// 如果发起方的钱不够了，就报错退出。
	if deposit < amount {
		panic("not enough money")
	}

	// 创建交易输入。
	for txIDString, indexes := range UTXOToPay {
		txID, err := hex.DecodeString(txIDString)
		if err != nil {
			panic(err)
		}
		for _, index := range indexes {
			newInputs = append(newInputs, transaction.NewTxi(txID, index, nil, wallet.Pubkey))
		}
	}

	// 创建交易输出。
	newOutputs = append(newOutputs, transaction.NewTxo(amount, to))

	// 如果需要找零，就多加一笔记录。
	if deposit > amount {
		newOutputs = append(newOutputs, transaction.NewTxo(deposit-amount, from))
	}

	// 将输入、输出存储进该次交易内。
	newTX := transaction.Transaction{
		ID:      nil,
		Inputs:  newInputs,
		Outputs: newOutputs,
	}
	newTX.ID = newTX.Hash()

	// 发起方对该次交易签名。
	c.SignTx(&newTX, wallet.Privkey)

	return &newTX
}
