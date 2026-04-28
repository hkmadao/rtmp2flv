package ext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/tcpserver"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

type FlvFileMediaInfoParam struct {
	IdCameraRecord string `json:"idCameraRecord"`
}

func ClientCameraRecordFileDuration(ctx *gin.Context) {
	idClient, ok := ctx.Params.Get("idClient")
	if !ok || idClient == "" {
		logs.Error("path param idClient is rquired")
		http.Error(ctx.Writer, "path param idClient is rquired", http.StatusBadRequest)
		return
	}

	clientInfo, err := base_service.ClientInfoSelectById(idClient)
	if err != nil {
		logs.Error("ClientInfoSelectById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

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

	// 缓冲一次回包，避免 HTTP 超时后 TCP 读协程卡在发送上。
	messageChan := make(chan *tcpserver.ResMessage, 1)
	done := make(chan struct{})
	defer close(done)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  clientInfo.ClientCode,
		MessageType: "flvFileMediaInfo",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
		Done:        done,
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
	sendReverseCommandErr := tcpserver.SendReverseCommand(clientInfo.Secret, &rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	if sendReverseCommandErr != nil {
		logs.Error("SendReverseCommand error: %v", err)
		if !sendReverseCommandErr.IsCustomError() {
			http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		result := common.ErrorResult(sendReverseCommandErr.Error())
		ctx.JSON(http.StatusOK, result)
		return
	}

	// 使用可停止的 timer，避免 NewTicker 留下不必要的计时资源。
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	select {
	case resMessage := <-messageChan:
		ctx.Data(http.StatusOK, gin.MIMEJSON, *resMessage.Data)
	case <-timer.C:
		logs.Error("read form client time out")
	case <-ctx.Request.Context().Done():
		logs.Error("client request canceled")
	}
}

type PlayParam struct {
	PlayerId       string `json:"playerId"`
	IdCameraRecord string `json:"idCameraRecord"`
	SeekSecond     uint64 `json:"seekSecond"`
}

func ClientCameraRecordFilePlay(ctx *gin.Context) {
	idClient, ok := ctx.Params.Get("idClient")
	if !ok || idClient == "" {
		logs.Error("path param idClient is rquired")
		http.Error(ctx.Writer, "path param idClient is rquired", http.StatusBadRequest)
		return
	}

	clientInfo, err := base_service.ClientInfoSelectById(idClient)
	if err != nil {
		logs.Error("ClientInfoSelectById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

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

	// 播放链路会持续传输分片，保留无缓冲通道做背压，通过 Done 处理取消。
	messageChan := make(chan *tcpserver.ResMessage)
	done := make(chan struct{})
	defer close(done)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  clientInfo.ClientCode,
		MessageType: "flvPlay",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
		Done:        done,
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
	sendReverseCommandErr := tcpserver.SendReverseCommand(clientInfo.Secret, &rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	if sendReverseCommandErr != nil {
		logs.Error("SendReverseCommand error: %v", err)
		if !sendReverseCommandErr.IsCustomError() {
			http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Error(ctx.Writer, sendReverseCommandErr.Error(), http.StatusInternalServerError)
		return
	}
	// 每收到一个分片就重置空闲计时；超时或 HTTP 取消会触发 TCP 读协程退出。
	idleTimer := time.NewTimer(1 * time.Minute)
	defer idleTimer.Stop()
Loop:
	for {
		select {
		case resMessage, ok := <-messageChan:
			if !ok {
				logs.Info("messageChan is closed, exit")
				break Loop
			}

			_, err := ctx.Writer.Write([]byte(*resMessage.Data))
			if err != nil {
				logs.Error("ctx write error: %v", err)
				break Loop
			}
			if !idleTimer.Stop() {
				select {
				case <-idleTimer.C:
				default:
				}
			}
			idleTimer.Reset(1 * time.Minute)
		case <-idleTimer.C:
			logs.Error("read form client time out")
			break Loop
		case <-ctx.Request.Context().Done():
			logs.Error("client request canceled")
			break Loop
		}
	}
}

type FetchMoreDataParam struct {
	PlayerId   string `json:"playerId"`
	SeekSecond uint64 `json:"seekSecond"`
}

func ClientCameraRecordFileFetch(ctx *gin.Context) {
	idClient, ok := ctx.Params.Get("idClient")
	if !ok || idClient == "" {
		logs.Error("path param idClient is rquired")
		http.Error(ctx.Writer, "path param idClient is rquired", http.StatusBadRequest)
		return
	}

	clientInfo, err := base_service.ClientInfoSelectById(idClient)
	if err != nil {
		logs.Error("ClientInfoSelectById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

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

	// 缓冲一次回包，避免客户端迟到响应导致 TCP 读协程滞留。
	messageChan := make(chan *tcpserver.ResMessage, 1)
	done := make(chan struct{})
	defer close(done)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  clientInfo.ClientCode,
		MessageType: "flvFetchMoreData",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
		Done:        done,
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
	sendReverseCommandErr := tcpserver.SendReverseCommand(clientInfo.Secret, &rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	if sendReverseCommandErr != nil {
		logs.Error("SendReverseCommand error: %v", err)
		if !sendReverseCommandErr.IsCustomError() {
			http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		result := common.ErrorResult(sendReverseCommandErr.Error())
		ctx.JSON(http.StatusOK, result)
		return
	}
	// 函数返回时停止 timer，避免频繁 fetch 时累计 ticker。
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	select {
	case resMessage := <-messageChan:
		ctx.Data(http.StatusOK, gin.MIMEJSON, *resMessage.Data)
	case <-timer.C:
		logs.Error("read form client time out")
	case <-ctx.Request.Context().Done():
		logs.Error("client request canceled")
	}
}

func ClientCameraAq(ctx *gin.Context) {
	idClient, ok := ctx.Params.Get("idClient")
	if !ok || idClient == "" {
		logs.Error("path param idClient is rquired")
		http.Error(ctx.Writer, "path param idClient is rquired", http.StatusBadRequest)
		return
	}

	clientInfo, err := base_service.ClientInfoSelectById(idClient)
	if err != nil {
		logs.Error("ClientInfoSelectById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	condition := common.AqCondition{}
	err = ctx.BindJSON(&condition)
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

	// 缓冲一次回包，避免摄像头查询取消后阻塞 TCP 读协程。
	messageChan := make(chan *tcpserver.ResMessage, 1)
	done := make(chan struct{})
	defer close(done)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  clientInfo.ClientCode,
		MessageType: "cameraAq",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
		Done:        done,
	}

	paramBytes, err := json.Marshal(condition)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sendReverseCommandErr := tcpserver.SendReverseCommand(clientInfo.Secret, &rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	if sendReverseCommandErr != nil {
		logs.Error("SendReverseCommand error: %v", sendReverseCommandErr)
		if !sendReverseCommandErr.IsCustomError() {
			http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		result := common.ErrorResult(sendReverseCommandErr.Error())
		ctx.JSON(http.StatusOK, result)
		return
	}

	// HTTP 请求取消时会关闭 Done，让 TCP 读协程退出。
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	select {
	case resMessage := <-messageChan:
		ctx.Data(http.StatusOK, gin.MIMEJSON, *resMessage.Data)
	case <-timer.C:
		logs.Error("read form client time out")
	case <-ctx.Request.Context().Done():
		logs.Error("client request canceled")
	}
}

func ClientCameraRecordAqPage(ctx *gin.Context) {
	idClient, ok := ctx.Params.Get("idClient")
	if !ok || idClient == "" {
		logs.Error("path param idClient is rquired")
		http.Error(ctx.Writer, "path param idClient is rquired", http.StatusBadRequest)
		return
	}

	clientInfo, err := base_service.ClientInfoSelectById(idClient)
	if err != nil {
		logs.Error("ClientInfoSelectById error: %v", err)
		result := common.ErrorResult("internal error")
		ctx.JSON(http.StatusOK, result)
		return
	}

	pageInfoInput := common.AqPageInfoInput{}
	err = ctx.BindJSON(&pageInfoInput)
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

	// 缓冲一次回包，避免历史录像查询取消后阻塞 TCP 读协程。
	messageChan := make(chan *tcpserver.ResMessage, 1)
	done := make(chan struct{})
	defer close(done)
	rcm := tcpserver.ReverseCommandMessage{
		ClientCode:  clientInfo.ClientCode,
		MessageType: "historyVideoPage",
		MessageId:   messageId,
		Created:     time.Now(),
		MessageChan: messageChan,
		Done:        done,
	}

	paramBytes, err := json.Marshal(pageInfoInput)
	if err != nil {
		logs.Error("param marshal error: %v", err)
		http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sendReverseCommandErr := tcpserver.SendReverseCommand(clientInfo.Secret, &rcm, string(paramBytes))
	defer tcpserver.ClearReverseCommand(messageId)
	if sendReverseCommandErr != nil {
		logs.Error("SendReverseCommand error: %v", err)
		if !sendReverseCommandErr.IsCustomError() {
			http.Error(ctx.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		result := common.ErrorResult(sendReverseCommandErr.Error())
		ctx.JSON(http.StatusOK, result)
		return
	}

	// HTTP 请求取消时会关闭 Done，让 TCP 读协程退出。
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	select {
	case resMessage := <-messageChan:
		ctx.Data(http.StatusOK, gin.MIMEJSON, *resMessage.Data)
	case <-timer.C:
		logs.Error("read form client time out")
	case <-ctx.Request.Context().Done():
		logs.Error("client request canceled")
	}
}
