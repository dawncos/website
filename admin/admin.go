package main

import (
	"time"

	_ "dawncos.com/admin/models"
	_ "dawncos.com/admin/routers"

	"dawncos.com/admin/utils"
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
)

func main() {
	utils.Che = cache.New(60*time.Minute, 120*time.Minute)
	beego.Run()
}
