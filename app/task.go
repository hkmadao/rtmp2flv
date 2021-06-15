package app

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/models"
)

func ClearToken() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("ClearToken panic : %v", r)
		}
	}()
	for {
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
	}
}

func DeleteExpiredAuthCode() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("DeleteExpiredAuthCode panic : %v", r)
		}
	}()
	for {
		css, err := models.CameraShareSelectAll()
		if err != nil {
			logs.Error("query camera error : %v", err)
		}
		for _, cs := range css {
			timeout := time.Now().After(cs.Created.Add(30 * 24 * time.Hour))
			if timeout {
				models.CameraShareDelete(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}
