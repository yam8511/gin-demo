package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	// _ "github.com/joho/godotenv/autoload" ----> uncomment will autoload .env file
)

func main() {
	setEnv()
	var err interface{}

	/// Create Gin Framework
	r := gin.Default()
	r.LoadHTMLGlob("view/*")
	r.GET("/ping", pong)
	r.NoRoute(NotFoundHandle)

	/// Home Demo
	r.Static("/assets", "./assets")
	r.Static("/images", "./images")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	/// Database Demo (Upper.io)
	upper := r.Group("/")
	{
		upper.GET("/upper", upperDemo)
	}

	/// WebSocket Setting
	wsserver := SocketInit()
	/// Socket Demo
	r.Static("/asset", "./asset")
	socket := r.Group("/", AccessAllowSetting)
	{
		socket.GET("/socket.io/", gin.WrapH(wsserver))
		socket.POST("/broadcast", BroadcastHandle(wsserver))
		socket.GET("/socket-demo", func(c *gin.Context) {
			if pusher, ok := c.Writer.(http.Pusher); ok {
				// Push is supported
				if err = pusher.Push("/asset/js/socket.io-1.3.7.js", nil); err != nil {
					log.Println("Server Push Error", err)
				}

				if err = pusher.Push("/asset/js/jquery-1.11.1.js", nil); err != nil {
					log.Println("Server Push Error", err)
				}
			} else {
				log.Println("Server Push is not supported!")
			}
			c.HTML(http.StatusOK, "chat.html", nil)
		})
	}

	/// GraphQL Demo
	todoInit()
	schemaInit()
	r.Static("/static", "./static")
	graphql := r.Group("/", AccessAllowSetting)
	{
		graphql.Any("/graphql", GraphQLHandle)
		graphql.GET("/graphiql", GraphIQLHandle)
		graphql.Any("/apollo-graphql", ApolloGraphQLHandle)
		graphql.GET("/apollo-todo", func(c *gin.Context) {
			c.HTML(http.StatusOK, "apollo-todo.html", nil)
		})
	}

	/// 宣告系統信號
	sigs := make(chan os.Signal, 1)
	exit := make(chan interface{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	/// 宣告設定 Server
	server := &http.Server{
		Addr:    PORT,
		Handler: r,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	go func() {
		/// 監聽伺服器
		log.Println("Server Listening on ", HOST+PORT)
		// go Bot.SendMessage(AdminChat, HOST+PORT+" 伺服器開啟了!", nil)
		err = server.ListenAndServe()
		// err = server.ListenAndServeTLS("server.crt", "server.key")

		// 如果監聽發生錯誤，通知系統人員
		if err != nil {
			if fmt.Sprintf("%v", err) != "http: Server closed" {
				log.Println("Server Error", err)
				NoticeSystemManager(err)
				close(exit)
			}
		}
	}()

	/// 設置 Ctrl + C 機制
	go func() {
		log.Println("結束程式請按 [Ctrl + C]")
		// 等待 Ctrl + C 的信號
		receivedSignel := <-sigs

		// 關閉伺服器
		err = server.Close()
		// 通知伺服器被關閉
		NoticeSystemManager(err)

		// 離開程式
		exit <- receivedSignel
	}()

	/// 等待 Ctrl + C 結束程式
	log.Printf("\nSignal: %v", <-exit)
	log.Println("程式結束")
}

// AccessAllowSetting : 伺服器存取權限設定
func AccessAllowSetting(c *gin.Context) {
	// AccessOrigin : 跨網域存取權限
	AccessOrigin := os.Getenv("ACCESS_ORIGIN")
	AccessCredentials := os.Getenv("ACCESS_CREDENTIAL")
	AccessHeaders := os.Getenv("ACCESS_HEADER")
	AccessMethods := os.Getenv("ACCESS_METHOD")

	c.Writer.Header().Set("Access-Control-Allow-Origin", AccessOrigin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", AccessCredentials)
	c.Writer.Header().Set("Access-Control-Allow-Headers", AccessHeaders)
	c.Writer.Header().Set("Access-Control-Allow-Methods", AccessMethods)
}

// NotFoundHandle : 404 Page
func NotFoundHandle(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
