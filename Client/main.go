package main

import (
	"Client/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

var serverAdd = "ws://123.172.7.3:8200/ws/join"

//默认启动
func main() {

	//log记录设置
	beego.SetLogger("file", `{"filename":"./logs/logs.log"}`)

	beego.SetLogFuncCall(true)

	go depth.DepthRun()

	client, _, err := websocket.DefaultDialer.Dial(serverAdd, nil)
	if err != nil {
		fmt.Println("websocket_join_error")
		return
	}

	fmt.Println("启动成功，有任何问题请加群联系管理员")
	//发送配置文件
	if config_data != "" {

		err = client.WriteMessage(websocket.TextMessage, []byte(config_data))
		if err != nil {
			fmt.Println("视频加速失败，请尝试关闭后重新启动")
		}
	}

	for {
		client.SetReadDeadline(time.Now().Add(60 * time.Second))
		_, message, err := client.ReadMessage()

		if err != nil {
			client.WriteMessage(websocket.TextMessage, []byte("websocket_ReadMessage_error"))
			fmt.Println("websocket_ReadMessage_error", err)
			go errLogs("websocket_ReadMessage_error")
			return
		}

		err = client.WriteMessage(websocket.TextMessage, content)
		if err != nil {
			fmt.Println("websocket_WriteMessage_error", err)
			go errLogs("websocket_WriteMessage_error")
			return
		}

	}

}
