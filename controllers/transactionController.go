package controllers

import (
	"blockchainDemo/blockchain_v2"
	"github.com/astaxie/beego"
)

type TransactionController struct {
	beego.Controller
}
var bc2 *blockchain_v2.BlockChain
func (this *TransactionController)Get() {
	// 创建区块链 不会让bc一直存储在服务器中
	bc2 = blockchain_v2.NewBlockChain("Alice")
	bc2.GetBalance("Alice") // 12.5
	//tx := blockchain_v2.NewTransaction("Alice", "Bob", 5,bc2) //
	//
	//coinbase := blockchain_v2.NewCoinbaseTx("Alice","")
	//bc2.AddBlock([]*blockchain_v2.Transaction{coinbase,tx})
	bc2.PrintChain()
	bc2.GetBalance("Alice")
	// 创世块
	this.Data["Genesis"] = bc2.Blocks[0]
	// 其他区块
	this.Data["blocks"] = bc2.Blocks[1:]
	this.TplName = "transaction.html"
}

// 添加了新区块
func (this *TransactionController) Post() {
	// TODO
	// 创建coinbase交易
	miner := this.GetString("miner")
	data := this.GetString("data")
	coinbase := blockchain_v2.NewCoinbaseTx(miner, data)
	// 创建一笔交易
	from := this.GetString("from")
	to := this.GetString("to")
	amount, err := this.GetFloat("amount")
	blockchain_v2.CheckErr("Get amount", err)
	tx := blockchain_v2.NewTransaction(from,to,amount,bc2)
	if tx == nil {
		beego.Info("Money is not enough!")
		// 创世块
		this.Data["Genesis"] = bc2.Blocks[0]
		// 其他区块
		this.Data["blocks"] = bc2.Blocks[1:]
		this.TplName = "transaction.html"
		return
	}
	bc2.AddBlock([]*blockchain_v2.Transaction{coinbase, tx})
	//addData := this.GetString("addData")
	//txs := blockchain_v2.NewTransaction(from, to, amount,bc2)
	//bc2.AddBlock(txs)
	// 每次添加新区块就打印信息
	bc2.PrintChain()
	// 创世块
	this.Data["Genesis"] = bc2.Blocks[0]
	// 其他区块
	this.Data["Blocks"] = bc2.Blocks[1:]
	this.TplName = "transaction.html"
}