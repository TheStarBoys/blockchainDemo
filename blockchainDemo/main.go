package main

import (
	_ "blockchainDemo/routers"
	"github.com/astaxie/beego"
)

func main() {
	bc := NewBlockChain()
	bc.AddBlock("Hello World")
	bc.PrintChain()
	beego.Run()
}

