package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
	message := HOST + PORT + " 伺服器關閉了!  "
	if err != nil {
		message += fmt.Sprintf("(%v)", err)
	} else {
		message += "(手動關閉)"
	}
	Bot.SendMessage(AdminChat, message, nil)
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
