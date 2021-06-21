package app

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/models"
)

type task struct {
}

func NewTask() *task {
	t := &task{}
	go t.clearToken()
	go t.deleteExpiredAuthCode()
	go t.offlineCamera()
	return t
}

func (t *task) clearToken() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
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

func (t *task) deleteExpiredAuthCode() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		css, err := models.CameraShareSelectAll()
		if err != nil {
			logs.Error("query CameraShare error : %v", err)
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

func (t *task) offlineCamera() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		css, err := models.CameraSelectAll()
		if err != nil {
			logs.Error("query camera error : %v", err)
		}
		exist := false
		for _, cs := range css {
			if cs.OnlineStatus != 1 {
				continue
			}
			rms.Range(func(key, value interface{}) bool {
				code := key.(string)
				if code == cs.Code {
					// r := value.(*RtmpManager)
					// if r.stop {
					// 	return true
					// }
					exist = true
				}
				return true
			})
			if !exist {
				cs.OnlineStatus = 0
				models.CameraUpdate(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}
