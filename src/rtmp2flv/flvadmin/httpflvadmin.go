package flvadmin

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvadmin/httpflvmanage"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/tcpserver"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dao/entity"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

var hfas *HttpFlvAdmin

type HttpFlvAdmin struct {
	hfms sync.Map
}

func init() {
	hfas = &HttpFlvAdmin{}
}

func GetSingleHttpFlvAdmin() *HttpFlvAdmin {
	return hfas
}

func (hfa *HttpFlvAdmin) AddHttpFlvManager(
	pktStream <-chan av.Packet,
	code string,
	codecs []av.CodecData,
) {
	hfm := httpflvmanage.NewHttpFlvManager(pktStream, code, codecs)
	hfa.hfms.Store(code, hfm)
	go func() {
		<-hfm.GetDone()
		hfa.hfms.Delete(code)
	}()
}

func (hfa *HttpFlvAdmin) StopWrite(code string) {
	v, ok := hfa.hfms.Load(code)
	if ok {
		ffw := v.(*httpflvmanage.HttpFlvManager)
		ffw.StopWrite()
	}
}

func (hfa *HttpFlvAdmin) StartWrite(code string) {
	v, ok := hfa.hfms.Load(code)
	if ok {
		ffw := v.(*httpflvmanage.HttpFlvManager)
		ffw.StopWrite()
		hfa.AddHttpFlvManager(ffw.GetPktStream(), code, ffw.GetCodecs())
	}
}

type RtmpPushParam struct {
	CameraCode string `json:"cameraCode"`
}

// 添加播放者
func (hfa *HttpFlvAdmin) AddHttpFlvPlayer(
	playerDone <-chan int,
	pulseInterval time.Duration,
	camera entity.Camera,
	writer http.ResponseWriter,
) (<-chan int, *common.Rtmp2FlvCustomError) {
	if !camera.OnlineStatus {
		return nil, common.CustomError("camera offline")
	}
	if !camera.Live {
		return nil, common.CustomError("camera live disabled")
	}
	v, b := hfa.hfms.Load(camera.Code)
	if b {
		hfm := v.(*httpflvmanage.HttpFlvManager)
		flvPlayerDone, err := hfm.AddHttpFlvPlayer(playerDone, pulseInterval, writer)
		if err != nil {
			return flvPlayerDone, common.InternalError(err)
		}
		if camera.FgPassive {
			go checkStopPushRtmp(flvPlayerDone, hfm, camera)
		}
		return flvPlayerDone, nil
	} else if camera.FgPassive {
		logs.Info("camera: %s rtmp push mode is passive, send command", camera.Code)
		messageId, err := utils.GenerateId()
		if err != nil {
			return nil, common.InternalError(err)
		}
		param := RtmpPushParam{CameraCode: camera.Code}
		messageChan := make(chan *tcpserver.ResMessage)
		rcm := tcpserver.ReverseCommandMessage{
			ClientCode:  camera.ClientInfo.ClientCode,
			MessageType: "startPushRtmp",
			MessageId:   messageId,
			Created:     time.Now(),
			MessageChan: messageChan,
		}

		paramBytes, err := json.Marshal(param)
		if err != nil {
			return nil, common.InternalError(err)
		}
		sendReverseCommandErr := tcpserver.SendReverseCommand(camera.ClientInfo.Secret, rcm, string(paramBytes))
		defer tcpserver.ClearReverseCommand(messageId)
		if sendReverseCommandErr != nil {
			logs.Error("SendReverseCommand error: %v", sendReverseCommandErr)
			if !sendReverseCommandErr.IsCustomError() {
				return nil, sendReverseCommandErr
			}

			return nil, sendReverseCommandErr
		}

		select {
		case resMessage := <-messageChan:
			result := common.AppResult{}
			err := json.Unmarshal(*resMessage.Data, &result)
			if err != nil {
				return nil, common.InternalError(err)
			}
			if result.IsFailed() {
				return nil, common.CustomError("client start push rtmp failed")
			}
			count := 0
			for {
				<-time.NewTicker(1 * time.Second).C
				count++
				v, b := hfa.hfms.Load(camera.Code)
				if b {
					hfm := v.(*httpflvmanage.HttpFlvManager)
					flvPlayerDone, err := hfm.AddHttpFlvPlayer(playerDone, pulseInterval, writer)
					if err != nil {
						return flvPlayerDone, common.InternalError(err)
					}
					go checkStopPushRtmp(flvPlayerDone, hfm, camera)
					return flvPlayerDone, nil
				}
				if count >= 59 {
					return nil, common.CustomError("client start push rtmp success, but the server not found rtmp connection")
				}
			}
		case <-time.NewTicker(1 * time.Minute).C:
			return nil, common.CustomError("read form client time out")
		}
	}

	return nil, common.CustomError("camera no connection")
}

// check exists player, stop client rtmp push
func checkStopPushRtmp(flvPlayerDone <-chan int, hfm *httpflvmanage.HttpFlvManager, camera entity.Camera) {
	<-flvPlayerDone
	checkStop(camera.Code, hfm)
}

func checkStop(cameraCode string, hfm *httpflvmanage.HttpFlvManager) {
	// first check
	existsPlayer := hfm.IsCameraExistsPlayer()
	if !existsPlayer {
		<-time.NewTicker(1 * time.Minute).C
		// sencod check
		existsPlayer = hfm.IsCameraExistsPlayer()
		if !existsPlayer {
			conditon := common.GetEqualCondition("code", cameraCode)
			camera, err := base_service.CameraFindOneByCondition(conditon)
			if err != nil {
				logs.Error("camera query error : %v", err)
				return
			}
			clientInfo, err := base_service.ClientInfoSelectById(camera.IdClientInfo)
			if err != nil {
				logs.Error("ClientInfo query error : %v", err)
				return
			}
			camera.ClientInfo = clientInfo
			messageId, err := utils.GenerateId()
			if err != nil {
				logs.Error("TickerCheckStopRtmp: %v", err)
			}
			messageChan := make(chan *tcpserver.ResMessage)
			rcm := tcpserver.ReverseCommandMessage{
				ClientCode:  camera.ClientInfo.ClientCode,
				MessageType: "stopPushRtmp",
				MessageId:   messageId,
				Created:     time.Now(),
				MessageChan: messageChan,
			}
			param := RtmpPushParam{CameraCode: camera.Code}
			paramBytes, err := json.Marshal(param)
			if err != nil {
				logs.Error("TickerCheckStopRtmp: %v", err)
			}
			sendReverseCommandErr := tcpserver.SendReverseCommand(camera.ClientInfo.Secret, rcm, string(paramBytes))
			defer tcpserver.ClearReverseCommand(messageId)
			if sendReverseCommandErr != nil {
				logs.Error("TickerCheckStopRtmp: %v", sendReverseCommandErr)
				return
			}
			select {
			case resMessage := <-messageChan:
				result := common.AppResult{}
				err := json.Unmarshal(*resMessage.Data, &result)
				if err != nil {
					logs.Error("TickerCheckStopRtmp: %v", err)
					return
				}
				if result.IsFailed() {
					logs.Error("TickerCheckStopRtmp client stop push rtmp failed")
					return
				}
			case <-time.NewTicker(1 * time.Minute).C:
				logs.Error("TickerCheckStopRtmp read form client time out")
				return
			}
		}
	}
}

// 更新sps、pps等信息
func (hfa *HttpFlvAdmin) UpdateCodecs(code string, codecs []av.CodecData) {
	rfw, ok := hfa.hfms.Load(code)
	if ok {
		rfw := rfw.(*httpflvmanage.HttpFlvManager)
		rfw.SetCodecs(codecs)
	}
}

func (hfa *HttpFlvAdmin) TickerCheckStopRtmp() {
	condition := common.GetEqualConditions([]common.EqualFilter{{Name: "onlineStatus", Value: true}, {Name: "fgPassive", Value: true}})
	css, err := base_service.CameraFindCollectionByCondition(condition)
	if err != nil {
		logs.Error("query camera error : %v", err)
	}
	for _, cs := range css {
		value, ok := hfa.hfms.Load(cs.Code)
		if !ok {
			continue
		}
		hfm := value.(*httpflvmanage.HttpFlvManager)
		go checkStop(cs.Code, hfm)
	}
}
