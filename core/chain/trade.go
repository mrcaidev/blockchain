package chain

import (
	tx "blockchain/transaction"
	"blockchain/utils"
	"blockchain/wallet"
	"encoding/hex"
	"fmt"
)

// 挖矿的奖励。
const subsidy = 10

// 创建一笔coinbase交易。
func NewCoinbaseTX(to string, data string) *tx.Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// 没有输入，只有一个输出。
	txi := tx.NewTXI([]byte{}, -1, nil, []byte(data))
	txo := tx.NewTXO(subsidy, to)
	tx := tx.Transaction{
		ID:      nil,
		Inputs:  []*tx.TXInput{txi},
		Outputs: []*tx.TXOutput{txo},
	}
	tx.ID = tx.Hash()
	return &tx
}

// 创建一笔UTXO交易。
func (chain *Chain) NewUTXOTX(from string, to string, cost int) *tx.Transaction {
	var (
		inputs  []*tx.TXInput
		outputs []*tx.TXOutput
	)

	wallets := wallet.LoadWallets()
	wallet := wallets.GetWallet(from)
	pubkeyHash := utils.GetPubkeyHash(wallet.Pubkey)

	deposit, UTXOToPay := chain.FindUTXOToPay(pubkeyHash, cost)

	// 如果发起方的钱不够了，就报错退出。
	if deposit < cost {
		panic("not enough money")
	}

	// 创建输入。
	for txid, utxo := range UTXOToPay {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			panic(err)
		}
		for _, out := range utxo {
			inputs = append(inputs, tx.NewTXI(txID, out, nil, wallet.Pubkey))
		}
	}

	// 创建输出。
	outputs = append(outputs, tx.NewTXO(cost, to))

	// 如果需要找零，就多加一笔记录。
	if deposit > cost {
		outputs = append(outputs, tx.NewTXO(deposit-cost, from))
	}

	// 将输入、输出存储进该次交易内。
	newTX := tx.Transaction{
		ID:      nil,
		Inputs:  inputs,
		Outputs: outputs,
	}
	newTX.ID = newTX.Hash()
	chain.SignTX(&newTX, wallet.Privkey)
	return &newTX
}
