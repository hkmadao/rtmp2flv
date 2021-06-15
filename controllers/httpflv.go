package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/result"
	"github.com/hkmadao/rtmp2flv/services"
)

func HttpFlvPlay(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.TrimSuffix(strings.TrimLeft(c.Request.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 3 || uris[0] != "live" {
		http.Error(c.Writer, "invalid path", http.StatusBadRequest)
		return
	}
	method := uris[1]
	code := uris[2]
	authCode := uris[3]
	r := result.Result{
		Code: 1,
		Msg:  "",
	}
	q := models.Camera{Code: code}
	camera, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("camera query error : %v", err)
		r.Code = 0
		r.Msg = "camera query error"
		c.JSON(http.StatusBadRequest, r)
		return
	}
	if !(method == "temp" || method == "permanent") {
		logs.Error("method error : %s", method)
		r.Code = 0
		r.Msg = "method error"
		c.JSON(http.StatusBadRequest, r)
		return
	}
	if method == "temp" {
		csq := models.CameraShare{CameraId: camera.Id, AuthCode: authCode}
		cs, err := models.CameraShareSelectOne(csq)
		if err != nil {
			logs.Error("CameraShareSelectOne error : %v", err)
			r.Code = 0
			r.Msg = "system error"
			c.JSON(http.StatusBadRequest, r)
			return
		}
		if time.Now().After(cs.Created.Add(7 * 24 * time.Hour)) {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			r.Code = 0
			r.Msg = "authCode expired"
			c.JSON(http.StatusBadRequest, r)
			return
		}

	}
	if method == "permanent" && authCode != camera.PlayAuthCode {
		logs.Error("AuthCodePermanent error : %s", authCode)
		r.Code = 0
		r.Msg = "authCode error"
		c.JSON(http.StatusBadRequest, r)
		return
	}
	//管理员可以主动中断播放
	done := make(chan interface{})
	heartChan := services.AddHttpFlvPlayer(done, code, c.Writer)
Loop:
	for {
		select {
		case heart := <-heartChan:
			if heart == 1 {
				continue
			}
			break Loop
		case <-time.After(10 * time.Second):
			logs.Info("player [%s] session [%s] timeout exit", code, c.Request.RemoteAddr)
			break Loop
		}
	}
	logs.Info("player [%s] session [%s] exit", code, c.Request.RemoteAddr)
}
