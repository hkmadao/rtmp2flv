package ext

import (
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvadmin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dao/entity"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
	ext_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/ext"
)

func CameraEnabled(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := base_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.Enabled = q.Enabled
	if q.Enabled != true {
		camera.OnlineStatus = false
	}
	_, err = base_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("enabled camera status %d error : %v", camera.Enabled, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	if q.Enabled != true {
		logs.Debug("close camera conn: %s", camera.Code)
		select {
		case codeStream <- camera.Code:
		case <-time.After(1 * time.Second):
		}
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraSaveVideoChange(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := base_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.SaveVideo = q.SaveVideo
	_, err = base_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("SaveVideo camera status %d error : %v", camera.SaveVideo, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	switch {
	case q.SaveVideo != true:
		logs.Info("camera [%s] stop save video", q.Code)
		flvadmin.GetSingleFileFlvAdmin().StopWrite(q.Code)
	case q.SaveVideo == true:
		flvadmin.GetSingleFileFlvAdmin().StartWrite(q.Code)
		logs.Info("camera [%s] start save video", q.Code)
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraLiveChange(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := base_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	camera.Live = q.Live
	_, err = base_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("Live camera status %d error : %v", camera.Live, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	if q.Live {
		flvadmin.GetSingleHttpFlvAdmin().StartWrite(q.Code)
	} else {
		flvadmin.GetSingleHttpFlvAdmin().StopWrite(q.Code)
	}

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

func CameraPlayAuthCodeReset(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	q := entity.Camera{}
	err := ctx.BindJSON(&q)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	camera, err := base_service.CameraSelectById(q.Id)
	if err != nil {
		logs.Error("query camera error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}
	playAuthCode := utils.GenarateRandName()
	camera.PlayAuthCode = playAuthCode
	_, err = base_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("Camera: %s PlayAuthCode reset error : %v", camera.Code, err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	flvadmin.GetSingleHttpFlvAdmin().StopWrite(q.Code)
	flvadmin.GetSingleHttpFlvAdmin().StartWrite(q.Code)

	result := common.SuccessResultWithMsg("succss", camera)
	ctx.JSON(http.StatusOK, result)
}

var codeStream = make(chan string)

func CodeStream() <-chan string {
	return codeStream
}

func CameraGetRecordFiles(ctx *gin.Context) {
	fileInfoList, err := ext_service.CameraFindRecordFiles()
	if err != nil {
		logs.Error("CameraGetRecordFiles error : %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	result := common.SuccessResultWithMsg("success", fileInfoList)
	ctx.JSON(http.StatusOK, result)
}
