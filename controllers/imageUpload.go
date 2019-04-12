package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"os"
)

type ImageUploadController struct {
	beego.Controller
}

// 文章更新 数据校验  路由 /api/article/update
func (this *ImageUploadController) ImageUpload() {
	// 获取文件信息
	file, information, err := this.GetFile("editormd-image-file")
	if err != nil {
		beego.Info(err)
	}
	file.Close()
	// 文件名
	fileName := information.Filename
	fmt.Println("上传图片名称", fileName)

	// 日期字符串
	//dateStr := beego.Date(time.Now(), "Ymd")
	// 创建文件夹
	filePath := "./static/uploads/"                       //+ dateStr
	if err1 := os.MkdirAll(filePath, 0777); err1 != nil { // 创建数据库目录
		panic("failed" + err.Error())
	}
	// folderPath := p.CreateDateDir("/static/uploads/")
	// 移动文件到创建好的文件夹内
	err2 := this.SaveToFile("editormd-image-file", filePath+"/"+fileName)
	if err2 == nil { // 如果没错，返回url,success,message
		this.Data["json"] = map[string]interface{}{
			"url":     "/static/uploads/" + "/" + fileName,
			"success": 1,
			"message": "upload success!",
		}
		this.ServeJSON()
	}
}
