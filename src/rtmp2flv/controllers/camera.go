package controllers

import (
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvadmin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/result"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
)

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
	page := result.Page{Total: len(cameras), Page: cameras}
	r.Data = page
	c.JSON(http.StatusOK, r)
}

func CameraDetail(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	cameraId := c.Query("id")
	if cameraId == "" {
		logs.Error("no cameraId found")
		r.Code = 0
		r.Msg = "no cameraId found"
		c.JSON(http.StatusOK, r)
		return
	}
	camera, err := models.CameraSelectById(cameraId)
	if err != nil {
		logs.Error("no camera found : %v", err)
		r.Code = 0
		r.Msg = "no camera found"
		c.JSON(http.StatusOK, r)
		return
	}
	r.Data = camera
	c.JSON(http.StatusOK, r)
}

func CameraEdit(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{
		Code: 1,
		Msg:  "",
	}
	q := models.Camera{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	if q.Id == "" || len(q.Id) == 0 {
		id, _ := utils.UUID()
		count, err := models.CameraCountByCode(q.Code)
		if err != nil {
			logs.Error("check camera is exist error : %v", err)
			r.Code = 0
			r.Msg = "check camera is exist"
			c.JSON(http.StatusOK, r)
			return
		}
		if count > 0 {
			logs.Error("camera code is exist error : %v", err)
			r.Code = 0
			r.Msg = "camera code is exist"
			c.JSON(http.StatusOK, r)
			return
		}
		q.Id = id
		q.Created = time.Now()
		playAuthCode, _ := utils.UUID()
		q.PlayAuthCode = playAuthCode
		_, err = models.CameraInsert(q)
		if err != nil {
			logs.Error("camera insert error : %v", err)
			r.Code = 0
			r.Msg = "camera insert error"
			c.JSON(http.StatusOK, r)
			return
		}
		c.JSON(http.StatusOK, r)
		return
	}
	count, err := models.CameraCountByCode(q.Code)
	if err != nil {
		logs.Error("check camera is exist error : %v", err)
		r.Code = 0
		r.Msg = "check camera is exist"
		c.JSON(http.StatusOK, r)
		return
	}
	if count > 1 {
		logs.Error("camera code is exist error : %v", err)
		r.Code = 0
		r.Msg = "camera code is exist"
		c.JSON(http.StatusOK, r)
		return
	}
	camera, _ := models.CameraSelectById(q.Id)
	camera.Code = q.Code
	camera.RtmpAuthCode = q.RtmpAuthCode
	// camera.Enabled = q.Enabled
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("camera insert error : %v", err)
		r.Code = 0
		r.Msg = "camera insert error"
		c.JSON(http.StatusOK, r)
		return
	}
	c.JSON(http.StatusOK, r)
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
	//close camera conn
	select {
	case codeStream <- camera.Code:
	case <-time.After(1 * time.Second):
	}

	c.JSON(http.StatusOK, r)
}

func CameraEnabled(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	q := models.Camera{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	camera, err := models.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		r.Code = 0
		r.Msg = "query camera error"
		c.JSON(http.StatusOK, r)
		return
	}
	camera.Enabled = q.Enabled
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("enabled camera status %d error : %v", camera.Enabled, err)
		r.Code = 0
		r.Msg = "enabled camera status %d error"
		c.JSON(http.StatusOK, r)
		return
	}
	if q.Enabled != 1 {
		//close camera conn
		select {
		case codeStream <- camera.Code:
		case <-time.After(1 * time.Second):
		}
	}

	c.JSON(http.StatusOK, r)
}

func CameraSaveVideoChange(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	q := models.Camera{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	camera, err := models.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		r.Code = 0
		r.Msg = "query camera error"
		c.JSON(http.StatusOK, r)
		return
	}
	camera.SaveVideo = q.SaveVideo
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("SaveVideo camera status %d error : %v", camera.SaveVideo, err)
		r.Code = 0
		r.Msg = "SaveVideo camera status %d error"
		c.JSON(http.StatusOK, r)
		return
	}
	switch {
	case q.SaveVideo != 1:
		logs.Info("camera [%s] stop save video", q.Code)
		flvadmin.GetSingleFileFlvAdmin().StopWrite(q.Code)
	case q.SaveVideo == 1:
		flvadmin.GetSingleFileFlvAdmin().StartWrite(q.Code)
		logs.Info("camera [%s] start save video", q.Code)
	}

	c.JSON(http.StatusOK, r)
}

func CameraLiveChange(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	q := models.Camera{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	camera, err := models.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		r.Code = 0
		r.Msg = "query camera error"
		c.JSON(http.StatusOK, r)
		return
	}
	camera.Live = q.Live
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("Live camera status %d error : %v", camera.Live, err)
		r.Code = 0
		r.Msg = "Live camera status %d error"
		c.JSON(http.StatusOK, r)
		return
	}
	switch {
	case q.Live != 1:
		flvadmin.GetSingleHttpflvAdmin().StopWrite(q.Code)
	case q.Live == 1:
		flvadmin.GetSingleHttpflvAdmin().StartWrite(q.Code)
	}

	c.JSON(http.StatusOK, r)
}

func CameraPlayAuthCodeReset(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{Code: 1, Msg: ""}
	q := models.Camera{}
	err := c.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	camera, err := models.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		r.Code = 0
		r.Msg = "query camera error"
		c.JSON(http.StatusOK, r)
		return
	}
	playAuthCode, _ := utils.UUID()
	camera.PlayAuthCode = playAuthCode
	_, err = models.CameraUpdate(camera)
	if err != nil {
		logs.Error("PlayAuthCode camera status %d error : %v", camera.PlayAuthCode, err)
		r.Code = 0
		r.Msg = "PlayAuthCode camera status %d error"
		c.JSON(http.StatusOK, r)
		return
	}
	flvadmin.GetSingleHttpflvAdmin().StopWrite(q.Code)
	flvadmin.GetSingleHttpflvAdmin().StartWrite(q.Code)

	c.JSON(http.StatusOK, r)
}

var codeStream = make(chan string)

func CodeStream() <-chan string {
	return codeStream
}
