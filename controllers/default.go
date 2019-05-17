package controllers

import (
	"blockchainDemo/blockchain"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}
var bc *blockchain.BlockChain
func init() {
	bc = blockchain.NewBlockChain()
	//loc, err := time.LoadLocation("Local")
	//blockchain.CheckErr("default init",err)
	//time.ParseInLocation(time.UnixDate,,loc)
	//timestamp := bc.Blocks[0].TimeStamp
	//dateTime := time.Unix(timestamp,0)
	//beego.Info("当前日期为：",dateTime.Format(time.UnixDate))
	bc.PrintChain()
}

func (c *MainController) Get() {
	// 创世块
	c.Data["Genesis"] = bc.Blocks[0]
	// 其他区块
	c.Data["blocks"] = bc.Blocks[1:]
	c.TplName = "index.html"
}
// 添加了新区块
func (this *MainController) Post() {
	addData := this.GetString("addData")
	bc.AddBlock(addData)
	// 每次添加新区块就打印信息
	bc.PrintChain()
	// 创世块
	this.Data["Genesis"] = bc.Blocks[0]
	// 其他区块
	this.Data["Blocks"] = bc.Blocks[1:]
	this.TplName = "index.html"
}