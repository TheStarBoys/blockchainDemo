package main

import (
	_ "github.com/TheStarBoys/blockchainDemo/controllers"
	_ "github.com/TheStarBoys/blockchainDemo/routers"
	"github.com/astaxie/beego"
	"time"
)

func main() {
	beego.AddFuncMap("AddOne", AddOne)
	beego.AddFuncMap("TimeStamp2Time", TimeStamp2Time)
	beego.Run()
}

// 让区块的序号加1， 得到正确序号
func AddOne(input int) (output int) {
	input++
	output = input
	return
}

// 由时间戳得到时间字符串
func TimeStamp2Time(input int64) (output string) {
	timestamp := input
	t := time.Unix(0, timestamp)
	//beego.Info("当前日期为：",t.Format(time.UnixDate))
	output = t.Format(time.UnixDate)
	return
}
