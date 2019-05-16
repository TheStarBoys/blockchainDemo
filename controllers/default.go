package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "index.html"
}

func (this *MainController) Post() {
	addData := this.GetString("addData")
	beego.Info("--------data: ",addData)
	this.Data["addData"] = addData
	this.TplName = "index.html"
}