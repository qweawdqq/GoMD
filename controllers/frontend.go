package controllers

import (
	"GoMD/models"
	"GoMD/tools"
	"fmt"
	"github.com/astaxie/beego"
)

type FrontendController struct {
	beego.Controller
}

var theme = models.GetOneConfig("Theme")
var layout = theme+"/layout.html"

//首页
func (this *FrontendController) Index() {
	page, _ := this.GetInt64("page")
	pageSize := tools.StringToInt64(models.GetOneConfig("PageSize"))
	if page < 1 {
		page = 1
	}
	list,count := models.LimitArticleDisplay(page,pageSize)
	this.Data["paging"] = tools.CreatePaging(page, pageSize, count)
	this.Data["list"] = list
	this.Data["config"] = models.ConfigList()
	this.Data["page"] = page
	this.Layout = layout
	this.TplName = theme+"/index.html"
}

//文章查看页面
func (this *FrontendController) Page() {
	id := this.GetString("id")
	if id == "" {
		id = "1"
	}
	article,_ := models.GetOneArticle(id)
	temp := *article
	label := models.SearchArticleCategory(temp[0].Cid)
	//这里加了打印看看都是啥
	//fmt.Println(article)
	this.Data["article"] = article
	this.Data["label"] = label
	this.Data["config"] = models.ConfigList()
	this.Layout = layout
	fmt.Println("风格数字是",theme)

	this.TplName = theme+"/page.html"
}

func (this *FrontendController) Search() {
	//查询页面
	keywords := this.GetString("keywords")

	var list *[]models.Article
	//list = models.GetAllArticle()
	////批量添加所有文章给引擎
	//for _, l := range *list{
	//	AddContent(&l)
	//}
	if isUseSearch {
		list = SearchSomgThing(keywords)
	}else {
	    list,_ = models.Search(keywords)
	}


	this.Data["list"] = list
	this.Data["config"] = models.ConfigList()
	this.Layout = layout
	this.TplName = theme+"/search.html"
}

func (this *FrontendController) About() {
	//关于 介绍页面
	this.Data["config"] = models.ConfigList()
	this.Layout = layout
	this.TplName = theme+"/about.html"
}

func (this *FrontendController) Archive() {
	//文章归档页面
	this.Data["config"] = models.ConfigList()
	data,_ := models.AllArticleList()
	this.Data["list"] = data
	this.Layout = layout
	this.TplName = theme+"/archive.html"
}