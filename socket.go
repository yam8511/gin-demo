package main

import (
	"fmt"
	"log"

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
