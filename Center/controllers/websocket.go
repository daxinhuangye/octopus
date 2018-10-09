package controllers

import (
	"Center/models"
	"fmt"
)

type WebSocketController struct {
	BaseController
}

//全局ws管理器
var Hub = models.NewHub()

func (this *WebSocketController) Get() {
	// Safe check.
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
	this.Display("websocket")
}

func (this *WebSocketController) Count() {

	echo := fmt.Sprintf("在线用户：%d", Hub.Count())
	this.Ctx.WriteString(echo)
}

var guestId = 100000

func (this *WebSocketController) Join() {

	guestId++
	nick_name := fmt.Sprintf("游客%d", guestId)
	channel := this.GetString("match")
	models.ServeWs(nick_name, channel, Hub, this.Ctx.ResponseWriter, this.Ctx.Request)
	fmt.Println("连接成功")
	this.Code = 1
	this.TraceJson()
}
