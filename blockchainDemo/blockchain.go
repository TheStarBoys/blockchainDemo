package main

import "fmt"

const splitLine = "-----------------------------------------------------------------------"

type BlockChain struct {
	blocks []*Block
}
// new新的区块链
func NewBlockChain() *BlockChain{
	//bc := BlockChain{}
	//bc.blocks = append(bc.blocks, NewGenesisBlock())
	//return &bc
	block := NewGenesisBlock()
	return &BlockChain{blocks:[]*Block{block}}
}

// 添加区块到区块链
func (bc *BlockChain)AddBlock(data string) {
	// 最后一个区块的哈希, 就是前哈希
	prevHash := bc.blocks[len(bc.blocks) - 1].Hash
	block := NewBlock(data, prevHash)
	bc.blocks = append(bc.blocks, block)
}

// 打印区块链
func (bc *BlockChain)PrintChain() {
	// 遍历区块链
	for _, block := range bc.blocks {
		fmt.Println(splitLine)
		fmt.Printf("Version:			%d\n",block.Version)
		fmt.Printf("PrevBlockHash:	%x\n",block.PrevBlockHash)
		fmt.Printf("Hash:			%x\n",block.Hash)
		fmt.Printf("MerkleRoot:		%x\n",block.MerkleRoot)
		fmt.Printf("TimeStamp:		%d\n",block.TimeStamp)
		fmt.Printf("Bits:			%d\n",block.Bits)
		fmt.Printf("Nonce:			%d\n",block.Nonce)
		fmt.Printf("Data:			%s\n",block.Data)
		fmt.Println(splitLine)
	}
}