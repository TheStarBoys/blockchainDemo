package blockchain

import (
	"crypto/sha256"
	"time"
)

type Block struct {
	Version			int64		`json "version""`
	PrevBlockHash	[]byte		`json "prevBlock"`
	Hash 			[]byte		`json "hash"`
	MerkleRoot		[]byte		`json "merkleRoot"`
	TimeStamp		int64		`json "timeStamp"`
	Bits			int64		`json "bits"`
	Nonce			int64		`json "nonce"`
	Data 			[]byte		`json "data"`
}
// new新区块
func NewBlock(data string, prevHash []byte) *Block {
	// 目前MerkleRoot 由粗略的hash运算得到
	mr := sha256.Sum256([]byte(data))
	block := Block{
		Version:1,
		PrevBlockHash:prevHash,
		MerkleRoot:mr[:],
		TimeStamp:time.Now().UnixNano(),
		Bits:targetBits,
		Data:[]byte(data),
	}
	// Nonce, hash 通过挖矿找
	pow := NewPoW(&block)
	pow.Mining()
	return &block
}
// new创世快
func NewGenesisBlock() *Block{
	return NewBlock("Welcome to Blockchain Demo by TheStarBoys" /*+ time.Now().Format(time.UnixDate)*/,[]byte{})
}