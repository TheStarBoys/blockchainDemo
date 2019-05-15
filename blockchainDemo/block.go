package main

import "time"

type Block struct {
	Version			int64
	PrevBlockHash	[]byte
	Hash 			[]byte
	MerkleRoot		[]byte
	TimeStamp		int64
	Bits			int64
	Nonce			int64
	Data 			[]byte
}
// new新区块
func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:1,
		PrevBlockHash:prevHash,
		MerkleRoot:[]byte{},
		TimeStamp:time.Now().UnixNano(),
		Bits:targetBits,
		Data:[]byte(data),
	}
	// Nonce, hash 通过挖矿找 TODO
	pow := NewPoW(&block)
	pow.Mining()
	return &block
}
// new创世快
func NewGenesisBlock() *Block{
	return NewBlock("The blockchain demo by TheStarBoy " + time.Now().Format("2006-01-02 15:04:05"),[]byte{})
}