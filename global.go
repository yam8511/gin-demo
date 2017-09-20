package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tucnak/telebot"
	"upper.io/db.v3/mysql"
)

// HOST : Host Name
var HOST string

// PORT : Port Number
var PORT string

// BotToken : TeleBot Token
var BotToken string

// AdminChat : Admin Group
var AdminChat telebot.Chat

// Bot : Telegram Bot
var Bot *telebot.Bot

// DBConfig : 資料庫設定
var DBConfig = mysql.ConnectionURL{}

func setEnv() {
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

	// Database Env
	dbHOST := os.Getenv("DB_HOST")
	dbDATABASE := os.Getenv("DB_DATABASE")
	dbUSER := os.Getenv("DB_USER")
	dbPASSWORD := os.Getenv("DB_PASSWORD")
	if dbHOST != "" {
		DBConfig.Host = dbHOST
	}
	if dbDATABASE != "" {
		DBConfig.Database = dbDATABASE
	}
	if dbUSER != "" {
		DBConfig.User = dbUSER
	}
	if dbPASSWORD != "" {
		DBConfig.Password = dbPASSWORD
	}

	// Telegram Env
	BotToken = os.Getenv("BOT_TOKEN")
	cahtid, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	CheckErrFatal(err, "Error Type: BOT_TOKEN")
	ChatID := int64(cahtid)

	/// Telegram - 伺服器關掉時通知
	AdminChat = telebot.Chat{ID: ChatID}
	Bot, err = telebot.NewBot(BotToken)
	CheckErrFatal(err, "NewBot Error", err)

	/// Set Gin Mode
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}
