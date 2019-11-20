package controllers

type ApiMonitorController struct {
	BaseController
}

func (self *ApiMonitorController) List() {
	self.Data["pageTitle"] = "API文档"
	self.display()
}
