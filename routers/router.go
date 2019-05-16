package routers

import (
	"blockchainDemo/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/add-block", &controllers.)
}
