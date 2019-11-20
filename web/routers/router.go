package routers

import (
	"dawncos.com/web/controllers"
	"github.com/astaxie/beego"
	"github.com/beego/admin"
)

func init() {
	admin.Run()
    beego.Router("/", &controllers.MainController{})
	beego.Router("/application/download", &controllers.FileOptDownloadController{}, "get:Download")
	beego.Router("/upload", &controllers.UploadController{})
}
