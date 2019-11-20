package controllers

import (
	"strings"
	"time"

	"dawncos.com/admin/libs"
	"dawncos.com/admin/models"
	"dawncos.com/admin/utils"
	"github.com/patrickmn/go-cache"
	"strconv"
)

type UserController struct {
	BaseController
}

func (self *UserController) Edit() {
	self.Data["pageTitle"] = "资料修改"
	id := self.userId
	Admin, _ := models.AdminGetById(id)
	row := make(map[string]interface{})
	row["id"] = Admin.Id
	row["login_name"] = Admin.LoginName
	row["real_name"] = Admin.RealName
	row["phone"] = Admin.Phone
	row["email"] = Admin.Email
	self.Data["admin"] = row
	utils.Che.Set("uid"+strconv.Itoa(self.user.Id), nil, cache.DefaultExpiration)
	self.display()
}

func (self *UserController) AjaxSave() {
	Admin_id, _ := self.GetInt("id")
	Admin, _ := models.AdminGetById(Admin_id)
	//修改
	Admin.Id = Admin_id
	Admin.UpdateTime = time.Now().Unix()
	Admin.UpdateId = self.userId
	Admin.LoginName = strings.TrimSpace(self.GetString("login_name"))
	Admin.RealName = strings.TrimSpace(self.GetString("real_name"))
	Admin.Phone = strings.TrimSpace(self.GetString("phone"))
	Admin.Email = strings.TrimSpace(self.GetString("email"))

	resetPwd := self.GetString("reset_pwd")
	if resetPwd == "1" {
		pwdOld := strings.TrimSpace(self.GetString("password_old"))
		pwdOldMd5 := libs.Md5([]byte(pwdOld + Admin.Salt))
		if Admin.Password != pwdOldMd5 {
			self.ajaxMsg("旧密码错误", MsgErr)
		}

		pwdNew1 := strings.TrimSpace(self.GetString("password_new1"))
		pwdNew2 := strings.TrimSpace(self.GetString("password_new2"))

		if pwdNew1 != pwdNew2 {
			self.ajaxMsg("两次密码不一致", MsgErr)
		}

		pwd, salt := libs.Password(4, pwdNew1)
		Admin.Password = pwd
		Admin.Salt = salt
	}
	Admin.UpdateTime = time.Now().Unix()
	Admin.UpdateId = self.userId
	Admin.Status = 1

	if err := Admin.Update(); err != nil {
		self.ajaxMsg(err.Error(), MsgErr)
	}
	self.ajaxMsg("", MsgOk)
}
