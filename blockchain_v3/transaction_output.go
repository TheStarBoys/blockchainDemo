package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// 代表交易的一个输出
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// Lock 签名一个输出
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - addressChecksumLen]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey 检查是否提供的公钥哈希被用于锁定输出
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput 创建一个新的交易输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{Value: value}
	txo.Lock([]byte(address))

	return txo
}

// TXOutputs 是TXOutput的集合
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize 序列化交易输出的集合
func (outs TXOutputs) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// DeserializeOutputs 反序列化交易输出的集合
func DeserializeOutputs(data []byte) TXOutputs {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var outs TXOutputs
	err := decoder.Decode(&outs)
	if err != nil {
		log.Panic(err)
	}

	return outs
}