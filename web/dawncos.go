package main

import (
	_ "dawncos.com/web/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/files", "static/file")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/images", "static/img")
	beego.Run()
}

