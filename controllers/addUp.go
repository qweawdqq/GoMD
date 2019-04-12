package controllers

import (
	"github.com/astaxie/beego"
)

type AddUpController struct {
	beego.Controller
}

// 文章更新 数据校验  路由 /api/article/update
func (this *AddUpController) AddUp() {
	//this.Layout = layout
	//this.TplName = theme + "/tongji.html"
	this.TplName = theme + "/tongji.html"
}
