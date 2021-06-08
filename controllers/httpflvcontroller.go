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
	sessionId := time.Now().Format("20060102150405")
	fw := &services.HttpFlvWriter{
		SessionId: sessionId,
		Writer:    c.Writer,
	}
	uri := strings.TrimSuffix(strings.TrimLeft(c.Request.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 2 || uris[0] != "live" {
		http.Error(c.Writer, "invalid path", http.StatusBadRequest)
		return
	}
	services.Hms[uris[1]].Fws[sessionId] = fw
	done := make(chan interface{})
	services.Hms[uris[1]].Fws[sessionId].Done = done
	<-done
	logs.Error("session %s exit", sessionId)
}

func CameraList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	cameras, err := models.CameraSelectAll()
	if err != nil {
		logs.Error("no camera found : %v", err)
		r.Code = 0
		r.Msg = "no camera found"
		c.JSON(http.StatusOK, r)
		return
	}
	r.Data = cameras
	c.JSON(http.StatusOK, r)
}

func CameraEdit(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	cameras, err := models.CameraSelectAll()
	if err != nil {
		logs.Error("no camera found : %v", err)
		c.JSON(http.StatusOK, cameras)
		return
	}
	c.JSON(http.StatusOK, cameras)
}

func CameraDelete(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	id, b := c.Params.Get("id")
	if !b {
		r.Code = 0
		r.Msg = "id is null"
		c.JSON(http.StatusOK, r)
		return
	}
	camera := models.Camera{Id: id}
	_, err := models.CameraDelete(camera)

	if err != nil {
		logs.Error("delete camera error : %v", err)
		r.Code = 0
		r.Msg = "delete camera error"
		c.JSON(http.StatusOK, r)
		return
	}
	c.JSON(http.StatusOK, r)
}
