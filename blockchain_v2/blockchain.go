package blockchain_v2

import (
	"fmt"
	"github.com/astaxie/beego"
)

const (
	splitLine = "-----------------------------------------------------------------------"
	genesisInfo = "Welcome to Blockchain Demo by TheStarBoys"
)
type BlockChain struct {
	Blocks []*Block
}
// new新的区块链
func NewBlockChain(address string) *BlockChain{
	//bc := BlockChain{}
	//bc.Blocks = append(bc.Blocks, NewGenesisBlock())
	//return &bc
	coinbase := NewCoinbaseTx(address, genesisInfo)
	block := NewGenesisBlock(coinbase)
	return &BlockChain{Blocks: []*Block{block}}
}

// 添加区块到区块链
func (bc *BlockChain)AddBlock(txs []*Transaction) {
	// 最后一个区块的哈希, 就是前哈希
	prevHash := bc.Blocks[len(bc.Blocks) - 1].Hash
	block := NewBlock(txs, prevHash)
	bc.Blocks = append(bc.Blocks, block)
}
// 颠倒区块顺序的函数 ***使用时小心， 由于是传的指针，会改变blocks的真正值， 用完后请再次调用，恢复原样
func BlocksRevers(blocks []*Block) []*Block{
	for i := 0; i < len(blocks) / 2; i++ {
		blocks[i], blocks[len(blocks) - i - 1] = blocks[len(blocks) - i - 1], blocks[i]
	}
	return blocks
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
		//fmt.Printf("Data:			%s\n",block.Data)
		fmt.Println("~~~~~~~~ Transactions ~~~~~~~")
		for _, tx := range block.Transactions {
			fmt.Printf("TXID:		%x\n", tx.TXID)
			fmt.Println("********** TXInputs **********")
			for _, input := range tx.TXInputs {
				fmt.Printf("TXID: 	%x\n",input.TXID)
				fmt.Printf("Vout: 	%d\n",input.Vout)
				fmt.Printf("ScriptSig: 	%s\n",input.ScriptSig)
				fmt.Println("#########")
			}
			fmt.Println("********** TXOutputs **********")
			for i, output := range tx.TXOutputs {
				fmt.Printf("outputIndex: 	%d\n", i)
				fmt.Printf("Value: 			%f\n", output.Value)
				fmt.Printf("ScriptPubKey: 	%s\n", output.ScriptPubKey)
				fmt.Println("#########")
			}
		}
		//fmt.Printf("send value:		%v\n",block.Transactions[0].TXOutputs[0].Value)
		fmt.Println("~~~~~~~~ Transactions ~~~~~~~")
		fmt.Println(splitLine)
	}
}
//找到支付用的UTXO集合 ,返回对应的交易ID和Vout的map, 总金额
func (bc *BlockChain)FindSuitableUXTOs(address string, amount float64) (map[string][]int64, float64) {
	//返回指定地址能够支配的UTXO的交易集合
	txs := bc.FindUTXOTransactions(address)
	validUTXOs := make(map[string][]int64)
	var total float64
	FIND:
	for _, tx := range txs {
		for index, output := range tx.TXOutputs {
			if output.CanBeUnlockedPubKey(address) {
				if total < amount {
					total += output.Value
					validUTXOs[string(tx.TXID)] =
						append(validUTXOs[string(tx.TXID)], int64(index))
				}else {
					break FIND
				}
			}
		}
	}
	return validUTXOs, total
}

//返回指定地址能够支配的UTXO的交易集合（包含未花费的output的交易集合）
func (bc *BlockChain) FindUTXOTransactions(address string) []Transaction {
	//包含目标UTXO的交易集合
	var UTXOTransactions []Transaction
	//存储使用过的UTXO的集合,map[交易id] : int64
	//0x1111111 : 0, 1  都是给Alice的转账
	spentUTXO := make(map[string][]int64)
	// 区块链需要从后往前遍历， ****否则无法达到功能 遍历完需要恢复顺序，否则会出错
	blocks := BlocksRevers(bc.Blocks)
	// 遍历区块
	for _, block := range blocks {
		// 遍历交易
		for _, tx := range block.Transactions {
			// 遍历输入
			if !tx.IsCoinbase() {
				//目的:找到已经消耗过的UTXO，把它们放到一个集合里
				//需要两个字段表示使用过的UTXO， 交易ID， output的索引
				for _, input := range tx.TXInputs {
					spentUTXO[string(input.TXID)] =
						append(spentUTXO[string(input.TXID)], input.Vout)
				}
			}
			OUTPUTS:
			// 遍历输出
			for currIndex, output := range tx.TXOutputs {
				//检查当前的output是否已经被消耗，如果消耗过，那么就进行下一个output检验
				if spentUTXO[string(tx.TXID)] != nil {
					// 非空，代表当前交易里面有消耗过的UTXO
					indexs := spentUTXO[string(tx.TXID)]
					for _, index := range indexs {
						if int64(currIndex) == index {
							// 当前的索引和消耗的索引相比较 若相同，说明output被消耗
							// 直接跳过，进行下一个output判断
							continue OUTPUTS
						}
					}
				}
				//如果当前地址是这个UTXO的所有者，就满足条件
				//目的找到所有能够支配的UTXO
				if output.CanBeUnlockedPubKey(address) {
					UTXOTransactions = append(UTXOTransactions, *tx)
				}

			}
		}
	}
	BlocksRevers(bc.Blocks)
	return UTXOTransactions
}
//寻找指定地址能够使用的UTXO
func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	//存储使用过的UTXO的集合,map[交易id] : int64
	//0x1111111 : 0, 1  都是给Alice的转账
	spentUTXO := make(map[string][]int64)
	var UTXOs []TXOutput
	txs := bc.FindUTXOTransactions(address)
	// 遍历交易
	for _, tx := range txs {
		// 遍历输入
		if !tx.IsCoinbase() {
			//目的:找到已经消耗过的UTXO，把它们放到一个集合里
			//需要两个字段表示使用过的UTXO， 交易ID， output的索引
			for _, input := range tx.TXInputs {
				spentUTXO[string(input.TXID)] =
					append(spentUTXO[string(input.TXID)], input.Vout)
			}
		}
	OUTPUTS:
		// 遍历输出
		for currIndex, output := range tx.TXOutputs {
			//检查当前的output是否已经被消耗，如果消耗过，那么就进行下一个output检验
			if spentUTXO[string(tx.TXID)] != nil {
				// 非空，代表当前交易里面有消耗过的UTXO
				indexs := spentUTXO[string(tx.TXID)]
				for _, index := range indexs {
					if int64(currIndex) == index {
						// 当前的索引和消耗的索引相比较 若相同，说明output被消耗
						// 直接跳过，进行下一个output判断
						continue OUTPUTS
					}
				}
			}
			//如果当前地址是这个UTXO的所有者，就满足条件
			//目的找到所有能够支配的UTXO
			if output.CanBeUnlockedPubKey(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}
	return UTXOs
}
//获取余额
func (bc *BlockChain)GetBalance(address string) {
	//获取没有消费过的utxo集合
	utxos := bc.FindUTXO(address)
	var total float64 = 0	//总金额
	//遍历所有的utxo获取金额总数
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("The Balance of %s is %f\n",address,total)
}
// 进行转账
func (bc *BlockChain)Send(from, to string, amount float64) {
	tx := NewTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
	beego.Info("send successfully!")
}