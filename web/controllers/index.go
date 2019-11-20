package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type UploadController struct {
	MainController
}

type FileOptDownloadController struct {
	MainController
}

func (c *MainController) Get() {
	c.Data["Website"] = "dawncos.com"
	c.Data["Email"] = "wxy8469163@gmail.com"
	c.Data["Download"] = "http://localhost:8080/application/download"
	c.Data["Upload"] = "http://localhost:8080/upload"
	c.TplName = "index.tpl"
}

func (u *UploadController) Get() {
	u.TplName="Upload.html"
}

func (u *UploadController) Post() {
	file, head, err := u.GetFile("file")
	if err != nil {
		u.Ctx.WriteString("获取文件失败")
		return
	}
	defer file.Close()

	filename := head.Filename
	err = u.SaveToFile("file","static/file/upload/" + filename)
	if err != nil {
		u.Ctx.WriteString("上传失败！")
	}else {
		u.Ctx.WriteString("上传成功！")
	}
}

func (f *FileOptDownloadController) Download() {
	//第一个参数是文件的地址，第二个参数是下载显示的文件的名称
	f.Ctx.Output.Download("static/file/app/app.jpg","root.jpg")
}
