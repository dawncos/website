package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"dawncos.com/admin/libs"
	"dawncos.com/admin/utils"

	"dawncos.com/admin/models"
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
)

const (
	MsgOk  = 0
	MsgErr = -1
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	user           *models.Admin
	userId         int
	userName       string
	loginName      string
	pageSize       int
	allowUrl       string
}

//前期准备
func (bc *BaseController) Prepare() {
	bc.pageSize = 20
	controllerName, actionName := bc.GetControllerAndAction()
	bc.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	bc.actionName = strings.ToLower(actionName)
	bc.Data["version"] = beego.AppConfig.String("version")
	bc.Data["siteName"] = beego.AppConfig.String("site.name")
	bc.Data["curRoute"] = bc.controllerName + "." + bc.actionName
	bc.Data["curController"] = bc.controllerName
	bc.Data["curAction"] = bc.actionName
	fmt.Println(bc.controllerName)
	if (strings.Compare(bc.controllerName, "index")) != 0 {
		bc.auth()
	}

	bc.Data["loginUserId"] = bc.userId
	bc.Data["loginUserName"] = bc.userName
}

//登录权限验证
func (bc *BaseController) auth() {
	arr := strings.Split(bc.Ctx.GetCookie("auth"), "|")
	bc.userId = 0
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		authId, _ := strconv.Atoi(idstr)
		if authId > 0 {
			var err error

			cheUser, found := utils.Che.Get("uid" + strconv.Itoa(authId))
			user := &models.Admin{}
			if found && cheUser != nil { //从缓存取用户
				user = cheUser.(*models.Admin)
			} else {
				user, err = models.AdminGetById(authId)
				utils.Che.Set("uid"+strconv.Itoa(authId), user, cache.DefaultExpiration)
			}
			if err == nil && password == libs.Md5([]byte(bc.getClientIp()+"|"+user.Password+user.Salt)) {
				bc.userId = user.Id

				bc.loginName = user.LoginName
				bc.userName = user.RealName
				bc.user = user
				bc.AdminAuth()
			}

			isHasAuth := strings.Contains(bc.allowUrl, bc.controllerName+"/"+bc.actionName)
			//不需要权限检查
			noAuth := "ajaxsave/ajaxdel/table/loginin/loginout/getnodes/start/show/ajaxapisave/index/group/public/env/code/apidetail"
			isNoAuth := strings.Contains(noAuth, bc.actionName)
			if isHasAuth == false && isNoAuth == false {
				bc.Ctx.WriteString("没有权限")
				bc.ajaxMsg("没有权限", MsgErr)
				return
			}
		}
	}

	if bc.userId == 0 && (bc.controllerName != "login" && bc.actionName != "loginin") {
		bc.redirect(beego.URLFor("LoginController.LoginIn"))
	}
}

func (bc *BaseController) AdminAuth() {
	cheMen, found := utils.Che.Get("menu" + strconv.Itoa(bc.user.Id))
	if found && cheMen != nil { //从缓存取菜单
		menu := cheMen.(*CheMenu)
		//fmt.Println("调用显示菜单")
		bc.Data["SideMenu1"] = menu.List1 //一级菜单
		bc.Data["SideMenu2"] = menu.List2 //二级菜单
		bc.allowUrl = menu.AllowUrl
	} else {
		// 左侧导航栏
		filters := make([]interface{}, 0)
		filters = append(filters, "status", 1)
		if bc.userId != 1 {
			//普通管理员
			adminAuthIds, _ := models.RoleAuthGetByIds(bc.user.RoleIds)
			adminAuthIdArr := strings.Split(adminAuthIds, ",")
			filters = append(filters, "id__in", adminAuthIdArr)
		}
		result, _ := models.AuthGetList(1, 1000, filters...)
		list := make([]map[string]interface{}, len(result))
		list2 := make([]map[string]interface{}, len(result))
		allow_url := ""
		i, j := 0, 0
		for _, v := range result {
			if v.AuthUrl != " " || v.AuthUrl != "/" {
				allow_url += v.AuthUrl
			}
			row := make(map[string]interface{})
			if v.Pid == 1 && v.IsShow == 1 {
				row["Id"] = int(v.Id)
				row["Sort"] = v.Sort
				row["AuthName"] = v.AuthName
				row["AuthUrl"] = v.AuthUrl
				row["Icon"] = v.Icon
				row["Pid"] = int(v.Pid)
				list[i] = row
				i++
			}
			if v.Pid != 1 && v.IsShow == 1 {
				row["Id"] = int(v.Id)
				row["Sort"] = v.Sort
				row["AuthName"] = v.AuthName
				row["AuthUrl"] = v.AuthUrl
				row["Icon"] = v.Icon
				row["Pid"] = int(v.Pid)
				list2[j] = row
				j++
			}
		}
		bc.Data["SideMenu1"] = list[:i]  //一级菜单
		bc.Data["SideMenu2"] = list2[:j] //二级菜单

		bc.allowUrl = allow_url + "/home/index"
		cheM := &CheMenu{}
		cheM.AllowUrl = bc.allowUrl
		cheM.List1 = bc.Data["SideMenu1"].([]map[string]interface{})
		cheM.List2 = bc.Data["SideMenu2"].([]map[string]interface{})
		utils.Che.Set("menu"+strconv.Itoa(bc.user.Id), cheM, cache.DefaultExpiration)
	}

}

type CheMenu struct {
	List1    []map[string]interface{}
	List2    []map[string]interface{}
	AllowUrl string
}

// 是否POST提交
func (bc *BaseController) isPost() bool {
	return bc.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (bc *BaseController) getClientIp() string {
	s := bc.Ctx.Request.RemoteAddr
	l := strings.LastIndex(s, ":")
	return s[0:l]
}

// 重定向
func (bc *BaseController) redirect(url string) {
	bc.Redirect(url, 302)
	bc.StopRun()
}

//加载模板
func (bc *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		tplname = bc.controllerName + "/" + bc.actionName + ".html"
	}
	bc.Layout = "public/layout.html"
	bc.TplName = tplname
}

//ajax返回
func (bc *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["status"] = msgno
	out["message"] = msg
	bc.Data["json"] = out
	bc.ServeJSON()
	bc.StopRun()
}

//ajax返回 列表
func (bc *BaseController) ajaxList(msg interface{}, msgno int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["msg"] = msg
	out["count"] = count
	out["data"] = data
	bc.Data["json"] = out
	bc.ServeJSON()
	bc.StopRun()
}

//分组公共方法
type groupList struct {
	Id        int
	GroupName string
}

func groupLists() (gl []groupList) {
	groupFilters := make([]interface{}, 0)
	groupFilters = append(groupFilters, "status", 1)
	groupResult, _ := models.GroupGetList(1, 1000, groupFilters...)
	for _, gv := range groupResult {
		groupRow := groupList{}
		groupRow.Id = int(gv.Id)
		groupRow.GroupName = gv.GroupName
		gl = append(gl, groupRow)
	}
	return gl
}

//获取单个分组信息
func getGroupInfo(gl []groupList, groupId int) (groupInfo groupList) {
	for _, v := range gl {
		if v.Id == groupId {
			groupInfo = v
		}
	}
	return
}

type sourceList struct {
	Id         int
	SourceName string
	GroupId    int
	GroupName  string
}

func sourceLists() (sl []sourceList) {

	grouplists := groupLists()
	var groupinfo groupList
	sourceFilters := make([]interface{}, 0)
	sourceFilters = append(sourceFilters, "status", 1)
	sourceResult, _ := models.ApiSourceGetList(1, 1000, sourceFilters...)
	for _, sv := range sourceResult {
		sourceRow := sourceList{}
		sourceRow.Id = int(sv.Id)
		sourceRow.GroupId = sv.GroupId
		groupinfo = getGroupInfo(grouplists, sv.GroupId)
		sourceRow.GroupName = groupinfo.GroupName
		sourceRow.SourceName = sv.SourceName
		sl = append(sl, sourceRow)
	}
	return sl
}

func getSourceInfo(gl []sourceList, sourceId int) (sourceInfo sourceList) {
	for _, v := range gl {
		if v.Id == sourceId {
			sourceInfo = v
		}
	}
	return
}

type envList struct {
	Id      int
	EnvName string
	EnvHost string
}

func envLists() (sl []envList) {
	envFilters := make([]interface{}, 0)
	envFilters = append(envFilters, "status__in", 1)
	envResult, _ := models.EnvGetList(1, 1000, envFilters...)
	for _, sv := range envResult {
		envRow := envList{}
		envRow.Id = int(sv.Id)
		envRow.EnvName = sv.EnvName
		envRow.EnvHost = sv.EnvHost
		sl = append(sl, envRow)
	}
	return sl
}

type templateList struct {
	Id           int
	TemplateName string
	Detail       string
}

func templateLists() (sl []templateList) {
	templateFilters := make([]interface{}, 0)
	templateFilters = append(templateFilters, "status", 1)
	templateResult, _ := models.TemplateGetList(1, 1000, templateFilters...)
	for _, sv := range templateResult {
		templateRow := templateList{}
		templateRow.Id = int(sv.Id)
		templateRow.TemplateName = sv.TemplateName
		templateRow.Detail = sv.Detail
		sl = append(sl, templateRow)
	}
	return sl
}

type codeList struct {
	Id     int
	Code   string
	Desc   string
	Detail string
}

func codeLists() (sl []codeList) {
	codeFilters := make([]interface{}, 0)
	codeFilters = append(codeFilters, "status", 1)
	codeResult, _ := models.CodeGetList(1, 1000, codeFilters...)
	for _, sv := range codeResult {
		codeRow := codeList{}
		codeRow.Id = int(sv.Id)
		codeRow.Code = sv.Code
		codeRow.Desc = sv.Desc
		codeRow.Detail = sv.Detail
		sl = append(sl, codeRow)
	}
	return sl
}

type apiPublicList struct {
	Id            int
	ApiPublicName string
	Sort          int
}

func apiPublicLists() (sl []apiPublicList) {
	apiPublicFilters := make([]interface{}, 0)
	apiPublicFilters = append(apiPublicFilters, "status", 1)
	apiPublicResult, _ := models.ApiPublicGetList(1, 1000, apiPublicFilters...)
	for _, sv := range apiPublicResult {
		apiPublicRow := apiPublicList{}
		apiPublicRow.Id = int(sv.Id)
		apiPublicRow.ApiPublicName = sv.ApiPublicName
		apiPublicRow.Sort = sv.Sort
		sl = append(sl, apiPublicRow)
	}
	return sl
}
