package transaction

import (
	"bytes"
	"encoding/gob"
)

// 交易输出集结构。
type TXOutputs struct {
	Outputs []*TXOutput // 交易输出集。
}

// 序列化交易输出集。
func (txos TXOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(txos)
	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

// 反序列化交易输出集。
func DeserializeTXOutputs(data []byte) TXOutputs {
	var txos TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&txos)
	if err != nil {
		panic(err)
	}

	return txos
}
