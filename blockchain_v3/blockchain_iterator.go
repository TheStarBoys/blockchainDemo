package main

import (
	"github.com/boltdb/bolt"
	"log"
)

// BlockchainIterator 用于迭代区块链的区块
type BlockchainIterator struct {
	currBlockHash  []byte
	db             *bolt.DB
}

// Next 返回下一个区块, 迭代的方向是从 最新块 --> 创世块
func (it *BlockchainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		currBlockData := b.Get(it.currBlockHash)
		block = DeserializeBlock(currBlockData)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	it.currBlockHash = block.PrevBlockHash

	return block
}
