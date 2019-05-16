package main

import (
	_ "blockchainDemo/controllers"
	_ "blockchainDemo/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("AddOne",AddOne)
	beego.Run()
}
// 让区块的序号加1， 得到正确序号
func AddOne(input int) (output int){
	//i, err := strconv.Atoi(input)
	//blockchain.CheckErr("AddOne",err)
	input++
	output = input
	return
}
