package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block 代表区块链中的一个区块
type Block struct {
	Timestamp     int64				// 时间戳
	Transactions  []*Transaction	// 交易集合
	PrevBlockHash []byte			// 前区块哈希
	Hash          []byte			// 当前区块的哈希
	Nonce         int 				// 随机值
	Height        int				// 当前块高度
}

// NewBlock 创建并返回区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		//Hash:,
		//Nonce:,
		Height:        height,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	return block
}

// NewGenesisBlock 创建并返回创世块
func NewGenesisBlock(coinbase *Transaction) *Block {
	if coinbase.IsCoinbase() == false {
		log.Panic("Transaction is not coinbase")
	}

	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// HashTransactions 返回区块中Merkle Root
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}

	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize 序列化区块，即将区块转换成字节码
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

// DeserializeBlock 反序列化一个区块, 即将字节码转换为区块
func DeserializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}