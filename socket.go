package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

// SocketInit : 初始化 Socket.io Server
func SocketInit() *socketio.Server {
	server, err := socketio.NewServer(nil)
	CheckErrFatal(err, err, "建立 WebSocket 錯誤")

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("melon")

		so.On("chat message", func(msg string) {
			server.BroadcastTo("melon", "chat message", msg)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		errMessage := fmt.Sprintf("WebSocket On Error\n (%v)\n", err)
		Bot.SendMessage(AdminChat, errMessage, nil)
	})

	return server
}

// BroadcastHandle : 處理廣播的請求
func BroadcastHandle(so *socketio.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		required := []string{}
		event, exists := c.GetQuery("event")
		if !exists {
			required = append(required, "event")
		}

		data, exists := c.GetQuery("data")
		if !exists {
			required = append(required, "data")
		}

		if len(required) > 0 {
			message := strings.Join(required, ",")
			c.JSON(http.StatusOK, gin.H{
				"error": message + " 為必填欄位!",
			})
			return
		}

		jsonData := map[string]interface{}{}
		var resData interface{}
		resData = data
		err := json.Unmarshal([]byte(data), &jsonData)
		if err == nil {
			resData = jsonData
		}

		so.BroadcastTo("melon", event, resData)
		c.JSON(http.StatusOK, gin.H{
			"event": event,
			"data":  resData,
		})
	}
}
