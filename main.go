package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
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

	/// Create Gin Framework
	r := gin.Default()
	v1 := r.Group("/", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", AccessOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	})
	{
		v1.GET("/ping", pong)
		v1.GET("/socket.io/", gin.WrapH(wsserver))
		v1.POST("/broadcast", BroadcastHandle(wsserver))
	}

	// Demo
	r.Static("/asset", "./asset")
	r.LoadHTMLGlob("asset/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

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

func authorizeDomain(host, port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host != host+port {
			c.String(http.StatusUnauthorized, "")
		}
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
