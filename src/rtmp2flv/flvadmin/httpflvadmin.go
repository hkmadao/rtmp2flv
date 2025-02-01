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
	v, b := hfa.hfms.Load(camera.Code)
	if b {
		hfm := v.(*httpflvmanage.HttpFlvManager)
		flvPlayerDone, err := hfm.AddHttpFlvPlayer(playerDone, pulseInterval, writer)
		return flvPlayerDone, common.InternalError(err)
	} else if camera.FgPassive {
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
					go checkStoppPushRtm(flvPlayerDone, hfm, camera)
					return flvPlayerDone, common.InternalError(err)
				}
				if count > 30 {
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
func checkStoppPushRtm(flvPlayerDone <-chan int, hfm *httpflvmanage.HttpFlvManager, camera entity.Camera) {
	<-flvPlayerDone
	// first check
	existsPlayer := hfm.IsCameraExistsPlayer()
	if !existsPlayer {
		<-time.NewTicker(1 * time.Minute).C
		// sencod check
		existsPlayer = hfm.IsCameraExistsPlayer()
		if !existsPlayer {
			messageId, err := utils.GenerateId()
			if err != nil {
				logs.Error("checkStoppPushRtm: %v", err)
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
				logs.Error("checkStoppPushRtm: %v", err)
			}
			sendReverseCommandErr := tcpserver.SendReverseCommand(camera.ClientInfo.Secret, rcm, string(paramBytes))
			defer tcpserver.ClearReverseCommand(messageId)
			if sendReverseCommandErr != nil {
				logs.Error("checkStoppPushRtm: %v", sendReverseCommandErr)
				return
			}
			select {
			case resMessage := <-messageChan:
				result := common.AppResult{}
				err := json.Unmarshal(*resMessage.Data, &result)
				if err != nil {
					logs.Error("checkStoppPushRtm: %v", err)
					return
				}
				if result.IsFailed() {
					logs.Error("checkStoppPushRtm client stop push rtmp failed")
					return
				}
			case <-time.NewTicker(1 * time.Minute).C:
				logs.Error("checkStoppPushRtm read form client time out")
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
