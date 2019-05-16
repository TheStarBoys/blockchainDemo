package blockchain

import "fmt"

const splitLine = "-----------------------------------------------------------------------"

type BlockChain struct {
	Blocks []*Block
}
// new新的区块链
func NewBlockChain() *BlockChain{
	//bc := BlockChain{}
	//bc.Blocks = append(bc.Blocks, NewGenesisBlock())
	//return &bc
	block := NewGenesisBlock()
	return &BlockChain{Blocks: []*Block{block}}
}

// 添加区块到区块链
func (bc *BlockChain)AddBlock(data string) {
	// 最后一个区块的哈希, 就是前哈希
	prevHash := bc.Blocks[len(bc.Blocks) - 1].Hash
	block := NewBlock(data, prevHash)
	bc.Blocks = append(bc.Blocks, block)
}

// 打印区块链
func (bc *BlockChain)PrintChain() {
	// 遍历区块链
	for _, block := range bc.Blocks {
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

// 将除GenesisBlock外的block的hash存在[][]byte中
//func (bc *BlockChain)SaveBlocksHash() (prevHashSlice, hashSlice []string){
//	for i, block := range bc.Blocks {
//		if i == 0 {	// 跳过Genesis block
//			continue
//		}
//		// 将区块的hash信息的十六进制字符串存到切片中
//		prevHashSlice = append(prevHashSlice, fmt.Sprintf("0x%x", block.PrevBlockHash))
//		hashSlice = append(hashSlice, fmt.Sprintf("0x%x", block.Hash))
//	}
//	return
//}