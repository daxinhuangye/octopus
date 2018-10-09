package main

import (
	"Center/controllers"

	_ "Center/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/beego/i18n"
)

/*
	LevelEmergency
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
*/

func main() {

	// Register template functions.
	beego.AddFuncMap("i18n", i18n.Tr)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	//log记录设置
	beego.SetLogger("file", `{"filename":"./logs/logs.log"}`)
	beego.SetLevel(beego.LevelDebug)
	if beego.AppConfig.String("runmode") == "prod" {
		beego.SetLevel(beego.LevelError)
	}
	beego.ErrorController(&controllers.ErrorController{})

	go controllers.Hub.Run()

	beego.Run()

}
