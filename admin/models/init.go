package models

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	urls := beego.AppConfig.String("mysqlurls")
	port := beego.AppConfig.String("mysqlport")
	user := beego.AppConfig.String("mysqluser")
	pass := beego.AppConfig.String("mysqlpass")
	name := beego.AppConfig.String("mysqldb")
	timezone := beego.AppConfig.String("mysqltimezone")
	if port == "" {
		port = "3306"
	}
	dsn := user + ":" + pass + "@tcp(" + urls + ":" + port + ")/" + name + "?charset=utf8"

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}

	_ = orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(Auth), new(Role), new(RoleAuth), new(Admin),
		new(Group), new(Env), new(Code), new(ApiSource), new(ApiDetail), new(ApiPublic), new(Template))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

func GetPrefix(name string) string {
	return beego.AppConfig.String("mysqlprefix") + name
}
