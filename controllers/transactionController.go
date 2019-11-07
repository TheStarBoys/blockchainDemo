package controllers

import (
	"github.com/TheStarBoys/blockchainDemo/blockchain_v2"
	"github.com/TheStarBoys/blockchainDemo/models"
	"github.com/astaxie/beego"
)

type TransactionController struct {
	beego.Controller
}
var bc2 *blockchain_v2.BlockChain
var accounts []*models.Account	// 需要传递指针， 容易错
func (this *TransactionController)Get() {
	// 创建区块链 不会让bc一直存储在服务器中
	bc2 = blockchain_v2.NewBlockChain("Alice")
	bc2.GetBalance("Alice") // 12.5
	//tx := blockchain_v2.NewTransaction("Alice", "Bob", 5,bc2) //
	//
	//coinbase := blockchain_v2.NewCoinbaseTx("Alice","")
	//bc2.AddBlock([]*blockchain_v2.Transaction{coinbase,tx})

	// 把所有账户集合放到"account" session里
	//var accounts []string
	//accounts = append(accounts, "Alice")
	//this.SetSession("accounts",accounts)
	//amounts := bc2.GetBalanceByAccounts(accounts)
	//for _, account := range accounts {
	//	balance := amounts[account]
	//	// 把accoutn : balance 存到session里
	//	this.SetSession(account, balance)
	//}
	//accounts = this.GetSession("accounts").([]string)
	// 清空切片
	accounts = accounts[:0]
	accounts = append(accounts, &models.Account{Name:"Alice",Amount:bc2.GetBalance("Alice")})
	bc2.PrintChain()
	bc2.GetBalance("Alice")
	// 创世块
	this.Data["Genesis"] = bc2.Blocks[0]
	// 其他区块
	this.Data["blocks"] = bc2.Blocks[1:]
	// accounts
	this.Data["accounts"] = accounts
	// 添加区块部分
	this.Data["money"] = 0.0
	this.Data["from"] = ""
	this.Data["to"] = ""
	this.Data["miner"] = ""
	this.Data["data"] = ""

	//this.Data["amounts"] = amounts
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
		this.Data["Blocks"] = bc2.Blocks[1:]
		this.Data["accounts"] = accounts
		// 添加区块部分
		this.Data["money"] = amount
		this.Data["from"] = from
		this.Data["to"] = to
		this.Data["miner"] = miner
		this.Data["data"] = data

		this.TplName = "transaction.html"
		return
	}
	if !checkAccountExist(from) {	// 如果没有该账户， 不能转账
	// TODO 如果没有该账户， 没有一个友好的提醒
		// 创世块
		// 创世块
		this.Data["Genesis"] = bc2.Blocks[0]
		// 其他区块
		this.Data["Blocks"] = bc2.Blocks[1:]
		this.Data["accounts"] = accounts
		// 添加区块部分
		this.Data["money"] = amount
		this.Data["from"] = from
		this.Data["to"] = to
		this.Data["miner"] = miner
		this.Data["data"] = data

		this.TplName = "transaction.html"
	}
	bc2.AddBlock([]*blockchain_v2.Transaction{coinbase, tx})
	// 一定要放在添加区块后，否则账户的余额不是最新的
	if !checkAccountExist(miner) {
		accounts = append(accounts, &models.Account{Name:miner,Amount:bc2.GetBalance(miner)})
	}else {
		// 存在就更新余额
		updateAmount(miner)
	}
	if !checkAccountExist(to) {
		accounts = append(accounts, &models.Account{Name:to,Amount:bc2.GetBalance(to)})
	}else {
		updateAmount(to)
	}
	if checkAccountExist(from) {	// 金钱来源存在， 才能更新金额
		updateAmount(from)
	}
	//addData := this.GetString("addData")
	//txs := blockchain_v2.NewTransaction(from, to, amount,bc2)
	//bc2.AddBlock(txs)
	// 每次添加新区块就打印信息
	bc2.PrintChain()
	// 创世块
	this.Data["Genesis"] = bc2.Blocks[0]
	// 其他区块
	this.Data["Blocks"] = bc2.Blocks[1:]
	this.Data["accounts"] = accounts
	// 添加区块部分
	this.Data["money"] = amount
	this.Data["from"] = from
	this.Data["to"] = to
	this.Data["miner"] = miner
	this.Data["data"] = data

	this.TplName = "transaction.html"
}
// 检查账户是否存在
func checkAccountExist(name string) bool {
	for _, account := range accounts {
		if account.Name == name {	// 如果账户存在
			return true
		}
	}
	return false
}
// 更新余额
func updateAmount(name string) {
	for _, account := range accounts {
		if account.Name == name {
			account.Amount = bc2.GetBalance(name)
		}
	}
}