package transaction

import (
	"bytes"
	"encoding/gob"
)

// 交易输出集结构。
type TxOutputs struct {
	List []*TxOutput // 交易输出列表。
}

// 序列化交易输出集。
func (txos *TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(txos)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

// 反序列化交易输出集。
func DeserializeTxOutputs(data []byte) *TxOutputs {
	var txos TxOutputs

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&txos)
	if err != nil {
		panic(err)
	}

	return &txos
}
