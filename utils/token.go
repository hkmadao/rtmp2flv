package utils

import (
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/google/uuid"
)

//start with date token
func NextToke() (string, error) {
	sessionId := time.Now().Format("20060102150405")
	id, err := uuid.NewRandom()
	if err != nil {
		logs.Error("Random error : %v", err)
		return "", err
	}
	return sessionId + "-" + id.String(), nil
}

//validate token
func TokenTimeOut(token string, duration time.Duration) bool {
	tokenTimeString := token[0:14]
	if len(tokenTimeString) != 14 {
		return true
	}
	tokenTime, err := time.Parse("20060102150405", tokenTimeString)
	if err != nil {
		return true
	}
	return time.Now().After(tokenTime.Add(duration))
}
