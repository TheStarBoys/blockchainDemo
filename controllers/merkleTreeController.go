package controllers

import (
	"github.com/TheStarBoys/blockchainDemo/blockchain_v3"
	"github.com/astaxie/beego"
)

type MerkleTreeController struct {
	beego.Controller
}

func (this *MerkleTreeController) Get() {
	this.TplName = "merkleTree.html"
}

func (this *MerkleTreeController) Post() {
	data1 := []byte(this.GetString("data1"))
	data2 := []byte(this.GetString("data2"))
	data3 := []byte(this.GetString("data3"))
	data4 := []byte(this.GetString("data4"))

	mTree := blockchain_v3.NewMerkleTree([][]byte{data1, data2, data3, data4})
	beego.Info("Start printTree")
	mTree.PrintTree()

	beego.Info("Start ToShow")
	this.Data["mNodes"] = mTree.ToShow()
	beego.Info("Start Tpl")
	this.TplName = "merkleTree.html"
}