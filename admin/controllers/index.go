package controllers

import (
	"fmt"
	"time"

	"dawncos.com/admin/models"
	"github.com/astaxie/beego"
)

type IndexController struct {
	BaseController
}

func (ic *IndexController) Index() {
	ic.Data["pageTitle"] = "道可思"
	ic.Data["ts"] = time.Now()

	gl := groupLists()
	ic.Data["grouplists"] = gl

	groupId, _ := ic.GetInt("id", 1)
	groupInfo, err := models.GroupGetById(groupId)

	if err != nil {
		fmt.Println("数据不存在")
	}

	//公共文档
	apiPublic, err := models.ApiPublicGetByIds(groupInfo.ApiPublicIds)
	ic.Data["apiPublic"] = apiPublic

	//环境
	// env, err := models.EnvGetByIds(groupInfo.EnvIds)
	// ic.Data["env"] = env

	// //状态码
	// code, err := models.CodeGetByIds(groupInfo.CodeIds)
	// ic.Data["code"] = code

	//接口
	apiMenu, _ := models.ApiTreeData(groupId)
	ic.Data["apiMenu"] = apiMenu
	ic.Data["groupId"] = groupId

	ic.TplName = "apidoc/index.html"
}

func (ic *IndexController) Public() {
	apiPublicId, _ := ic.GetInt("id", 1)
	apiPublic, _ := models.ApiPublicGetById(apiPublicId)
	ic.Data["apiPublic"] = apiPublic
	ic.TplName = "apidoc/apipublic.html"
}

func (ic *IndexController) Env() {
	groupId, _ := ic.GetInt("id", 0)
	groupInfo, _ := models.GroupGetById(groupId)
	env, _ := models.EnvGetByIds(groupInfo.EnvIds)
	ic.Data["env"] = env
	ic.TplName = "apidoc/env.html"
}

func (ic *IndexController) Code() {
	groupId, _ := ic.GetInt("id", 0)
	groupInfo, _ := models.GroupGetById(groupId)
	code, _ := models.CodeGetByIds(groupInfo.CodeIds)
	ic.Data["code"] = code
	ic.TplName = "apidoc/code.html"
}

func (ic *IndexController) ApiDetail() {
	id, _ := ic.GetInt("id", 0)
	detail, _ := models.ApiFullDetailById(id)
	row := make(map[string]interface{})
	row["id"] = detail.Id
	row["source_id"] = detail.SourceId
	row["api_url"] = detail.ApiUrl
	row["api_name"] = detail.ApiName
	row["detail"] = detail.Detail
	row["status"] = detail.Status
	row["create_name"] = detail.CreateName
	row["update_name"] = detail.UpdateName
	row["audit_name"] = detail.AuditName
	row["audit_status"] = AUDIT_STATUS[detail.Status]
	row["method"] = REQUEST_METHOD[detail.Method]
	row["audit_time"] = beego.Date(time.Unix(detail.AuditTime, 0), "Y-m-d H:i:s")
	row["update_time"] = beego.Date(time.Unix(detail.UpdateTime, 0), "Y-m-d H:i:s")

	ic.Data["pageTitle"] = "查看 " + detail.ApiName
	ic.Data["Detail"] = row
	ic.TplName = "apidoc/apidetail.html"
}
