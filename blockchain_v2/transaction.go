package blockchain_v2

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
)

const reward = 12.5
type Transaction struct {
	TXID		[]byte		// 交易ID
	TXInputs	[]TXInput	// 交易输入
	TXOutputs	[]TXOutput	// 交易输出
}

type TXInput struct {
	TXID		[]byte	// 交易ID
	Vout		int64	// 所引用的输出的索引
	ScriptSig	string	// 解锁脚本
}
// 能否和解锁脚本匹配
func (input *TXInput)CanUnlockUTXOSig(unlock string) bool {
	return input.ScriptSig == unlock
}

type TXOutput struct {
	Value			float64	// 转移金额
	ScriptPubKey	string	// 锁定脚本
}
// 能否和锁定脚本匹配
func (output *TXOutput)CanBeUnlockedPubKey(unlock string) bool {
	return output.ScriptPubKey == unlock
}
// 设置交易ID
func (tx *Transaction)SetTxId() {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	CheckErr("SetTxId", err)
	hash := sha256.Sum256(buf.Bytes())
	tx.TXID = hash[:]
}
// 创建coinbase交易
func NewCoinbaseTx(address,data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("reward to %s %f BTC", address, reward)
	}
	input := TXInput{TXID:[]byte{},Vout:-1,ScriptSig:data}
	output := TXOutput{Value:reward,ScriptPubKey:address}
	tx := Transaction{[]byte{},[]TXInput{input},[]TXOutput{output}}
	tx.SetTxId()
	return &tx
}
// 校验是否是coinbase
func (tx *Transaction)IsCoinbase() bool {
	if len(tx.TXInputs) == 1 {
		if len(tx.TXInputs[0].TXID) == 0 && tx.TXInputs[0].Vout == -1{
			return true
		}
	}
	return false
}
// 新增一笔交易
func NewTransaction(from, to string, amount float64, bc *BlockChain)  *Transaction {
	// key : value = txid : vout
	validUTXOs := make(map[string][]int64)
	var total float64
	validUTXOs, total = bc.FindSuitableUXTOs(from, amount)
	if total < amount {
		beego.Info("Not enough money")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	//1.创建inputs
	//进行output到input的转换
	//遍历有效UTXO的合集
	for txId, outputIndexes := range validUTXOs {
		//遍历所有引用的UTXO的索引，每一个索引需要创建一个input
		for _, index := range outputIndexes {
			input := TXInput{TXID:[]byte(txId),Vout:index,ScriptSig:from}
			inputs = append(inputs, input)
		}
	}
	//2.创建outputs
	//给对方转账
	output := TXOutput{amount,to}
	outputs = append(outputs, output)
	//找零
	if total > amount {
		output := TXOutput{total - amount, from}
		outputs = append(outputs, output)
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetTxId()
	return &tx
}