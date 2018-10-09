package routers

import (
	"Center/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// Register routers.
	beego.Router("/", &controllers.IndexController{})

	beego.Router("/join", &controllers.IndexController{}, "post:Join")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")
	beego.Router("/ws/count", &controllers.WebSocketController{}, "get:Count")
}
