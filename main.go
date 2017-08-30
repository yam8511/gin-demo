package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// _ "github.com/joho/godotenv/autoload" ----> uncomment will autoload .env file
	"github.com/tucnak/telebot"
)

// HOST : Host Name
var HOST string

// PORT : Port Number
var PORT string

// AccessOrigin : 跨網域存取權限
var AccessOrigin string

// BotToken : TeleBot Token
var BotToken string

// AdminChat : Admin Group
var AdminChat telebot.Chat

// Bot : Telegram Bot
var Bot *telebot.Bot

func main() {
	/// 設定參數
	envFile := flag.String("e", "", "指定 env 檔案名稱")
	flag.Parse()

	/// 讀取 ENV 設定檔
	var err interface{}
	if *envFile == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(*envFile)
	}
	CheckErrFatal(err, "讀取 .env 錯誤")

	/// 存取 Env 變數
	GinMode := os.Getenv("GIN_MODE")
	HOST = os.Getenv("HOST")
	PORT = ":" + os.Getenv("PORT")
	AccessOrigin = os.Getenv("ACCESS_ORIGIN")
	BotToken = os.Getenv("BOT_TOKEN")
	cahtid, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	CheckErrFatal(err, "Error Type: BOT_TOKEN")
	ChatID := int64(cahtid)

	/// Telegram - 伺服器關掉時通知
	AdminChat = telebot.Chat{ID: ChatID}
	Bot, err = telebot.NewBot(BotToken)
	CheckErrFatal(err, "NewBot Error", err)

	/// WebSocket Setting
	wsserver := SocketInit()

	/// Set Gin Mode
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	/// Create Gin Framework
	r := gin.Default()
	r.Static("/asset", "./asset")
	r.LoadHTMLGlob("asset/*")
	r.GET("/ping", pong)

	/// Socket Demo
	r.GET("/socket.io/", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", AccessOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}, gin.WrapH(wsserver))
	r.POST("/broadcast", BroadcastHandle(wsserver))
	r.GET("/socket-demo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	/// GraphQL Demo
	todoInit()
	r.Any("/graphql", GraphQLHandle)
	r.GET("/graphiql", GraphIQLHandle)

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
		err = server.ListenAndServe()
		Bot.SendMessage(AdminChat, HOST+PORT+" 伺服器開啟了!", nil)

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

// CheckErrFatal : 確認錯誤，如果有錯誤則結束程式
func CheckErrFatal(err interface{}, msg ...interface{}) {
	if err != nil {
		if len(msg) == 0 {
			log.Fatalln(err)
		}
		log.Fatalln(msg, err)
	}
}

// NoticeSystemManager : 通知系統人員
func NoticeSystemManager(err interface{}) {
	message := HOST + " 伺服器關閉了!  "
	if err != nil {
		message += fmt.Sprintf("(%v)", err)
	} else {
		message += "(手動關閉)"
	}
	Bot.SendMessage(AdminChat, message, nil)
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// AllInput : Request取所有輸入參數的方法
func AllInput(req *http.Request) map[string]interface{} {
	v := make(map[string]interface{})
	jsonData := make(map[string]interface{})

	err := req.ParseForm()
	if err == nil {
		for key, value := range req.Form {
			if len(value) == 1 {
				v[key] = value[0]
			} else {
				v[key] = value
			}
			fmt.Println("key: ", key, "value: ", value)
		}
	}

	err = json.NewDecoder(req.Body).Decode(&jsonData)
	if err == nil {
		for key, value := range jsonData {
			v[key] = value
		}
	}

	return v
}
