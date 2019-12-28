package main

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs 找到并返回未消费的交易输出
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k,v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if accumulated >= amount {
					return nil
				}
				if out.IsLockedWithKey(pubkeyHash) {
					accumulated += out.Value
					// TODO: Is this a bug? outIdx is not vout
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO 找到用公钥哈希锁定的UTXO
func (u *UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions 返回UTXO集中交易的数量
func (u *UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex 重新构建UTXO集
func (u *UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		b, err := tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		UTXO := u.Blockchain.FindUTXO()

		for txID, outputs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("reindex get key: %v, value: %v", key, outputs)

			err = b.Put(key, outputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// Update 通过区块来更新UTXO集合中的信息
// 该区块是最新块
func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, vin := range tx.Vin {
					unspentOuts := TXOutputs{}
					outsData := b.Get(vin.TXID)

					outs := DeserializeOutputs(outsData)
					for outIdx, out := range outs.Outputs {
						if outIdx == vin.Vout {
							continue
						}

						unspentOuts.Outputs = append(unspentOuts.Outputs, out)
					}

					if len(unspentOuts.Outputs) == 0 {
						err := b.Delete(vin.TXID)
						if err != nil {
							return err
						}
					} else {
						err := b.Put(vin.TXID, unspentOuts.Serialize())
						if err != nil {
							return err
						}
					}
				}
			}

			newOuts := TXOutputs{}
			for _, out := range tx.Vout {
				newOuts.Outputs = append(newOuts.Outputs, out)
			}

			err := b.Put(tx.ID, newOuts.Serialize())
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}