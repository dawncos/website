package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"dawncos.com/admin/libs"
	"dawncos.com/admin/models"
	"dawncos.com/admin/utils"
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
)

type LoginController struct {
	BaseController
}

//登录 TODO:XSRF过滤
func (lc *LoginController) LoginIn() {
	if lc.userId > 0 {
		lc.redirect(beego.URLFor("HomeController.Index"))
	}
	beego.ReadFromRequest(&lc.Controller)
	if lc.isPost() {

		username := strings.TrimSpace(lc.GetString("username"))
		password := strings.TrimSpace(lc.GetString("password"))

		if username != "" && password != "" {
			user, err := models.AdminGetByName(username)
			fmt.Println(user)
			flash := beego.NewFlash()
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == 0 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = lc.getClientIp()
				user.LastLogin = time.Now().Unix()
				user.Update()
				utils.Che.Set("uid"+strconv.Itoa(user.Id), user, cache.DefaultExpiration)
				authkey := libs.Md5([]byte(lc.getClientIp() + "|" + user.Password + user.Salt))
				lc.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)

				lc.redirect(beego.URLFor("HomeController.Index"))
			}
			flash.Error(errorMsg)
			flash.Store(&lc.Controller)
			lc.redirect(beego.URLFor("LoginController.LoginIn"))
		}
	}
	lc.TplName = "login/login.html"
}

//登出
func (lc *LoginController) LoginOut() {
	lc.Ctx.SetCookie("auth", "")
	lc.redirect(beego.URLFor("LoginController.LoginIn"))
}

func (lc *LoginController) NoAuth() {
	lc.Ctx.WriteString("没有权限")
}
