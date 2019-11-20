package routers

import (
	"dawncos.com/admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 默认登录
	beego.Router("/", &controllers.IndexController{}, "*:Index")
	beego.Router("/admin/login", &controllers.LoginController{}, "*:LoginIn")
	beego.Router("/admin/login_out", &controllers.LoginController{}, "*:LoginOut")
	beego.Router("/admin/no_auth", &controllers.LoginController{}, "*:NoAuth")

	beego.Router("/home", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/start", &controllers.HomeController{}, "*:Start")
	//beego.AutoRouter(&controllers.ApiController{})
	//beego.AutoRouter(&controllers.ApiSourceController{})
	//beego.AutoRouter(&controllers.ApiPublicController{})
	//beego.AutoRouter(&controllers.TemplateController{})
	//beego.AutoRouter(&controllers.IndexController{})
	//// beego.AutoRouter(&controllers.ApiMonitorController{})
	//beego.AutoRouter(&controllers.EnvController{})
	//beego.AutoRouter(&controllers.CodeController{})
	//
	//beego.AutoRouter(&controllers.GroupController{})
	//beego.AutoRouter(&controllers.AuthController{})
	//beego.AutoRouter(&controllers.RoleController{})
	//beego.AutoRouter(&controllers.AdminController{})
	//beego.AutoRouter(&controllers.UserController{})

}
