package controllers

import (
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/result"
	"github.com/hkmadao/rtmp2flv/utils"
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

func CameraEdit(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	r := result.Result{
		Code: 1,
		Msg:  "",
	}
	params := make(map[string]interface{})
	err := c.BindJSON(&params)
	if err != nil {
		logs.Error("param error : %v", err)
		r.Code = 0
		r.Msg = "param error"
		c.JSON(http.StatusOK, r)
		return
	}

	enabled, _ := params["enabled"].(int)
	q := models.Camera{
		Code:         params["code"].(string),
		RtmpAuthCode: params["rtmpAuthCode"].(string),
		Enabled:      enabled,
	}
	if params["id"] != nil {
		q.Id = params["id"].(string)
	}
	if q.Id == "" || len(q.Id) == 0 {
		id, _ := utils.NextToke()
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
		playAuthCode, _ := utils.NextToke()
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
	_, err = models.CameraUpdate(q)
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
	c.JSON(http.StatusOK, r)
}
