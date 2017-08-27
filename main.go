package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
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
	/// 讀取 ENV 設定檔
	err := godotenv.Load()
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
	AdminChat := &telebot.Chat{ID: ChatID}
	Bot, err = telebot.NewBot(BotToken)
	CheckErrFatal(err, "NewBot Error", err)
	defer func(BotToken string, ChatID int64) {
		err := recover()
		message := HOST + " 伺服器關閉了!  "
		if err != nil {
			message += fmt.Sprintf("(%v)", err)
		}
		Bot.SendMessage(AdminChat, message, nil)
	}(BotToken, ChatID)

	// WebSocket Setting
	wsserver := SocketInit()

	// // Set Gin Mode
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

	// Socket Demo
	r.GET("/socket.io/", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", AccessOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}, gin.WrapH(wsserver))
	r.POST("/broadcast", BroadcastHandle(wsserver))
	r.GET("/socket-demo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	// GraphQL Demo
	r.Any("/graphql", GraphQLHandle)
	r.GET("/graphiql", GraphIQLHandle)

	// Start Server
	r.Run(PORT)
}

// CheckErrFatal : 確認錯誤，如果有錯誤則結束程式
func CheckErrFatal(err interface{}, msg ...interface{}) {
	if err != nil {
		if len(msg) == 0 {
			log.Fatalln(err)
		}
		log.Fatalln(msg)
	}
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
