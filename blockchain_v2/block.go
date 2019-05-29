package blockchain_v2

import (
	"bytes"
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
	//Data 			[]byte		`json "data"`
	Transactions	[]*Transaction
}
// new新区块
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	// 目前MerkleRoot 由交易ID代替
	block := Block{
		Version:1,
		PrevBlockHash:prevHash,
		MerkleRoot:[]byte{},
		TimeStamp:time.Now().UnixNano(),
		Bits:targetBits,
		Transactions:txs}
	block.HashTransactions()
	// Nonce, hash 通过挖矿找
	pow := NewPoW(&block)
	pow.Mining()
	return &block
}
// new创世块
func NewGenesisBlock(coinbase *Transaction) *Block{
	genesis := NewBlock([]*Transaction{coinbase}, []byte{})
	genesis.HashTransactions()
	//return NewBlock("Welcome to Blockchain Demo by TheStarBoys" /*+ time.Now().Format(time.UnixDate)*/,[]byte{})
	return genesis
}
//粗略模拟merkle tree，将交易的哈希值进行拼接，生成root hash
func (block *Block)HashTransactions() []byte{
	var txHashes [][]byte
	txs := block.Transactions
	//遍历交易
	for _, tx := range txs {
		//[]byte
		txHashes = append(txHashes, tx.TXID)
	}
	//对二维切片进行拼接，生成一维切片
	data := bytes.Join(txHashes, []byte{})
	hash := sha256.Sum256(data)
	block.MerkleRoot = hash[:]
	return hash[:]
}