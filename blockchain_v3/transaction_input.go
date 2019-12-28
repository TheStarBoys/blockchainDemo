package main

import "bytes"

// 代表交易的一个输入
type TXInput struct {
	TXID      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UsesKey 检查输入使用了指定输入解锁一个输出
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(pubKeyHash, lockingHash) == 0
}