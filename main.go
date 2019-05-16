package main

import (
	_ "blockchainDemo/routers"
	"encoding/json"
	"github.com/astaxie/beego"
	"golang.org/x/exp/errors/fmt"
)

var bc *BlockChain

func main() {
	bc = NewBlockChain()
	bc.PrintChain()
	// 返回当前区块hash的函数
	beego.AddFuncMap("AddBlockData",AddBlockData)
	beego.AddFuncMap("LoadGenesisHash",LoadGenesisHash)
	beego.Run()
}
// 返回当前区块链的json数据
func AddBlockData (data string) (jsonStr string) {
	// 用一个[][]byte 来装hash格式化后的字符串
	// 每次添加区块， 都应该把区块信息存进去
	// 添加区块
	beego.Info("Add Block Data begin--------")
	bc.AddBlock(data)
	bc.PrintChain()
	j, err := json.MarshalIndent(bc,"","")
	CheckErr("AddBlockData", err)
	beego.Info(string(j))
	return string(j)
}
// 加载时，获得GenesisBlock的hash
func LoadGenesisHash() (hash string){
	// 格式化输出一个十六进制的hash字符串
	return fmt.Sprintf("0x%x", bc.Blocks[0].Hash)
}