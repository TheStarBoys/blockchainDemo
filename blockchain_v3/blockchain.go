package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain 实现与数据库的交互
type Blockchain struct {
	db  *bolt.DB
	tip []byte // 最后一个区块的哈希
}

// CreateBlockchain 创建一个新的区块链
func CreateBlockchain(address, nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	coinbase := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(coinbase)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}

		err = b.Put([]byte(genesis.Hash), genesis.Serialize())
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}

		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{db, tip}

	return &bc
}

// NewBlockchain 从已有的数据库中, 得到区块链数据, 新创建一个区块链对象
func NewBlockchain(nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{db, tip}

	return &bc
}

// AddBlock 保存区块到区块链中
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockInDb := b.Get(block.Hash)
		if blockInDb != nil {
			return fmt.Errorf("block already exist")
		}

		lastBlockHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastBlockHash)
		lastBlock := DeserializeBlock(lastBlockData)
		if lastBlock.Height < block.Height {
			return fmt.Errorf("block height error")
		}

		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}

		bc.tip = block.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// FindTransaction 通过交易ID找到对应交易
func (bc *Blockchain) FindTransaction(TXID []byte) (*Transaction, error) {
	it := bc.Iterator()

	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, TXID) == 0 {
				return tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil, fmt.Errorf("Transaction is not found")
}

// FindUTXO 找到所有未花费的交易输出并返回已移除已花费输出的交易集合
// TXID -> TXOutputs
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	it := bc.Iterator()

	for {
		block := it.Next()


		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, output := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOutIndx := range spentTXOs[txID] {
						if spentOutIndx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, output)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.TXID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator 返回一个区块链迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currBlockHash: bc.tip,
		db:            bc.db,
	}
}

// GetBestHeight 返回最新块的高度
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock *Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastBlockData := b.Get([]byte("l"))
		lastBlock = DeserializeBlock(lastBlockData)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlock 通过区块的哈希值找到一个区块, 并返回
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block *Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockData := b.Get(blockHash)
		if blockData == nil || len(blockData) == 0 {
			return fmt.Errorf("block is not found")
		}
		block = DeserializeBlock(blockData)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return *block, nil
}

// GetBlockHashes 返回区块链中所有区块哈希的集合
func (bc *Blockchain) GetBlockHashes() [][]byte {
	it := bc.Iterator()
	var hashes [][]byte

	for {
		block := it.Next()

		hashes = append(hashes, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return hashes
}

// MineBlock 根据提供的交易挖出一个新区块
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastBlockHash []byte
	var lastBlockHeight int

	for _, tx := range transactions {
		// TODO: ignore transaction if it's not valid
		if bc.VerifyTransaction(tx) == false {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		lastBlockHash = b.Get([]byte("l"))
		lastBlockData := b.Get(lastBlockHash)
		lastBlock := DeserializeBlock(lastBlockData)

		lastBlockHeight = lastBlock.Height

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastBlockHash, lastBlockHeight)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

// SignTransaction 签名交易的输入
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.TXID)
		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.ID)] = *prevTx
	}

	tx.Sign(privKey, prevTxs)
}

// VerifyTransaction 校验交易输入的签名
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.TXID)
		if err != nil {
			log.Panic(err)
		}

		prevTxs[hex.EncodeToString(prevTx.ID)] = *prevTx
	}

	return tx.Verify(prevTxs)
}

// dbExists 检测是否已经存在数据库
func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}