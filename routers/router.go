package routers

import (
	"blockchainDemo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// blockchain页面
    beego.Router("/", &controllers.MainController{})
	beego.Router("/blockchain", &controllers.MainController{})

    // block页面
    beego.Router("/block", &controllers.BlockController{})

    // transaction页面
    beego.Router("/transaction", &controllers.TransactionController{})
}
