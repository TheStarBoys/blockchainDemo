package controllers

import "github.com/astaxie/beego"

type AddBlockController struct {
	beego.Controller
}

func (this *AddBlockController)Get() {
	addData := this.GetString("addData")
	beego.Info("--------data: ",addData)
	this.TplName = "index.html"
}