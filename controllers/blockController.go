package controllers

import "github.com/astaxie/beego"

type BlockController struct {
	beego.Controller
}

func (this *BlockController)Get() {
	this.TplName = "block.html"
}