package ext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/tcpserver"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
)

type FlvFileMediaInfoParam struct {
	IdCameraRecord string `json:"idCameraRecord"`
}

func ClientCameraRecordFileDuration(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()

	idCameraRecord, ok := ctx.Params.Get("idCameraRecord")
	if !ok || idCameraRecord == "" {
		logs.Error("path param idCameraRecord is rquired")
		http.Error(ctx.Writer, "path param idCameraRecord is rquired", http.StatusBadRequest)
		return
	}

	messageId, err := utils.GenerateId()
	if err != nil {
		logs.Error("generate messageId error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  "demo",
		MessageType: "flvFileMediaInfo",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
	}
	param := FlvFileMediaInfoParam{
		IdCameraRecord: idCameraRecord,
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tcpserver.SendReverseCommand(rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	select {
	case messageBytes := <-messageChan:
		logs.Info(string(messageBytes))
		ctx.Data(http.StatusOK, gin.MIMEJSON, messageBytes)
	case <-time.NewTicker(1 * time.Minute).C:
		logs.Error("read form client time out")
	}
}

type PlayParam struct {
	PlayerId       string `json:"playerId"`
	IdCameraRecord string `json:"idCameraRecord"`
	SeekSecond     uint64 `json:"seekSecond"`
}

func ClientCameraRecordFilePlay(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	playerId := ctx.Query("playerId")
	if playerId == "" {
		logs.Error("query param playerId is rquired")
		http.Error(ctx.Writer, "query param playerId is rquired", http.StatusBadRequest)
		return
	}

	idCameraRecord, ok := ctx.Params.Get("idCameraRecord")
	if !ok || idCameraRecord == "" {
		logs.Error("path param idCameraRecord is rquired")
		http.Error(ctx.Writer, "path param idCameraRecord is rquired", http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if seekSecond == "" {
		logs.Error("query param seekSecond is rquired")
		http.Error(ctx.Writer, "query param seekSecond is rquired", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("query param seekSecond need uint")
		http.Error(ctx.Writer, "query param seekSecond need uint", http.StatusBadRequest)
		return
	}

	messageId, err := utils.GenerateId()
	if err != nil {
		logs.Error("generate messageId error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  "demo",
		MessageType: "flvPlay",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
	}
	param := PlayParam{
		IdCameraRecord: idCameraRecord,
		PlayerId:       playerId,
		SeekSecond:     seekSecondUint,
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tcpserver.SendReverseCommand(rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
Loop:
	for {
		select {
		case messageBytes, ok := <-messageChan:
			if !ok {
				logs.Info("messageChan is closed, exit")
				break Loop
			}

			_, err := ctx.Writer.Write(messageBytes)
			if err != nil {
				logs.Error("ctx write error: %v", err)
				break Loop
			}
		case <-time.NewTicker(1 * time.Minute).C:
			logs.Error("read form client time out")
			break Loop
		}
	}
}

type FetchMoreDataParam struct {
	PlayerId   string `json:"playerId"`
	SeekSecond uint64 `json:"seekSecond"`
}

func ClientCameraRecordFileFetch(ctx *gin.Context) {
	defer func() {
		if result := recover(); result != nil {
			logs.Error("system painc : %v \nstack : %v", result, string(debug.Stack()))
		}
	}()
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	playerId := ctx.Query("playerId")
	if playerId == "" {
		logs.Error("query param playerId failed")
		http.Error(ctx.Writer, "query param playerId is rquired", http.StatusBadRequest)
		return
	}

	seekSecond := ctx.Query("seekSecond")
	if playerId == "" {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "query param seekSecond is rquired", http.StatusBadRequest)
		return
	}
	seekSecondUint, err := strconv.ParseUint(seekSecond, 10, 64)
	if err != nil {
		logs.Error("get param seekSecond failed")
		http.Error(ctx.Writer, "query param seekSecond need uint", http.StatusBadRequest)
		return
	}

	messageId, err := utils.GenerateId()
	if err != nil {
		logs.Error("generate messageId error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  "demo",
		MessageType: "flvFetchMoreData",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
	}
	param := FetchMoreDataParam{
		PlayerId:   playerId,
		SeekSecond: seekSecondUint,
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tcpserver.SendReverseCommand(rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	select {
	case messageBytes := <-messageChan:
		logs.Info(string(messageBytes))
		ctx.Data(http.StatusOK, gin.MIMEJSON, messageBytes)
	case <-time.NewTicker(1 * time.Minute).C:
		logs.Error("read form client time out")
	}
}

func ClientCameraAq(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	condition := common.AqCondition{}
	err := ctx.BindJSON(&condition)
	if err != nil {
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}
	messageId, err := utils.GenerateId()
	if err != nil {
		logs.Error("generate messageId error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  "demo",
		MessageType: "cameraAq",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
	}

	paramBytes, err := json.Marshal(condition)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tcpserver.SendReverseCommand(rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	select {
	case messageBytes := <-messageChan:
		logs.Info(string(messageBytes))
		ctx.Data(http.StatusOK, gin.MIMEJSON, messageBytes)
	case <-time.NewTicker(1 * time.Minute).C:
		logs.Error("read form client time out")
	}
}

func ClientCameraRecordAqPage(ctx *gin.Context) {
	// ctx.Writeresult.Header().Set("Access-Control-Allow-Origin", "*")
	pageInfoInput := common.AqPageInfoInput{}
	err := ctx.BindJSON(&pageInfoInput)
	if err != nil {
		ctx.AbortWithError(500, err)
		logs.Error("param error : %v", err)
		result := common.ErrorResult(fmt.Sprintf("param error : %v", err))
		ctx.JSON(http.StatusOK, result)
		return
	}

	messageId, err := utils.GenerateId()
	if err != nil {
		logs.Error("generate messageId error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan []byte)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  "demo",
		MessageType: "historyVideoPage",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
	}

	paramBytes, err := json.Marshal(pageInfoInput)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tcpserver.SendReverseCommand(rcm, string(paramBytes))
	defer func() {
		tcpserver.ClearReverseCommand(messageId)
	}()
	select {
	case messageBytes := <-messageChan:
		logs.Info(string(messageBytes))
		ctx.Data(http.StatusOK, gin.MIMEJSON, messageBytes)
	case <-time.NewTicker(1 * time.Minute).C:
		logs.Error("read form client time out")
	}
}
