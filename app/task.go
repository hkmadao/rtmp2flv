package app

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/utils"
)

func ClearToken() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("ClearToken pain : %v", r)
		}
	}()
	deleteTokens := []string{}
	// 遍历所有sync.Map中的键值对
	tokens.Range(func(k, v interface{}) bool {
		if time.Now().After(v.(time.Time).Add(30 * time.Minute)) {
			deleteTokens = append(deleteTokens, k.(string))
		}
		return true
	})
	for _, v := range deleteTokens {
		tokens.Delete(v)
	}
	<-time.After(24 * time.Hour)
	ClearToken()
}

func UpdateAuthCode() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("UpdateAuthCode pain : %v", r)
		}
	}()
	cameras, err := models.CameraSelectAll()
	if err != nil {
		logs.Error("query camera error : %v", err)
	}
	for _, camera := range cameras {
		timeout := utils.TokenTimeOut(camera.AuthCodeTemp, 10080*time.Minute)
		if timeout {
			authCodeTemp, _ := utils.NextToke()
			camera.AuthCodeTemp = authCodeTemp
			models.CameraUpdate(camera)
		}
	}
	<-time.After(10 * time.Minute)
	UpdateAuthCode()
}
