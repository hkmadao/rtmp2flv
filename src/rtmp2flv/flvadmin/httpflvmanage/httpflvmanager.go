package httpflvmanage

import (
	"io"
	"runtime/debug"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvadmin/httpflvmanage/httpflvwriter"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/utils"
)

type HttpFlvManager struct {
	done      chan int
	pktStream <-chan av.Packet
	code      string
	codecs    []av.CodecData
	hfws      sync.Map
}

func (hfm *HttpFlvManager) GetCode() string {
	return hfm.code
}

func (hfm *HttpFlvManager) SetCodecs(codecs []av.CodecData) {
	hfm.codecs = codecs
	hfm.hfws.Range(func(key, value interface{}) bool {
		wi := value.(*httpflvwriter.HttpFlvWriter)
		wi.SetCodecs(hfm.codecs)
		return true
	})
}

func (hfm *HttpFlvManager) GetDone() <-chan int {
	return hfm.done
}

func (hfm *HttpFlvManager) GetPktStream() <-chan av.Packet {
	return hfm.pktStream
}

func (hfm *HttpFlvManager) GetCodecs() []av.CodecData {
	return hfm.codecs
}

func NewHttpFlvManager(pktStream <-chan av.Packet, code string, codecs []av.CodecData) *HttpFlvManager {
	hfm := &HttpFlvManager{
		done:      make(chan int),
		pktStream: pktStream,
		code:      code,
		codecs:    codecs,
	}
	camera, err := models.CameraSelectOne(models.Camera{Code: code})
	if err != nil {
		logs.Error("query camera error : %v", err)
		return hfm
	}
	if camera.OnlineStatus != 1 {
		return hfm
	}
	if camera.Live != 1 {
		go func() {
			for {
				select {
				case <-hfm.GetDone():
					return
				case _, ok := <-hfm.pktStream:
					if !ok {
						return
					}
				}
			}
		}()
		return hfm
	}
	go hfm.flvWrite()
	return hfm
}

func (hfm *HttpFlvManager) StopWrite() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
			}
		}()
		close(hfm.done)
	}()
}

//Write extends to writer.Writer
func (hfm *HttpFlvManager) flvWrite() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for pkt := range utils.OrDonePacket(hfm.done, hfm.pktStream) {
		hfm.hfws.Range(func(key, value interface{}) bool {
			wi := value.(*httpflvwriter.HttpFlvWriter)
			select {
			case wi.GetPktStream() <- pkt:
				// logs.Debug("flvWrite pkt")
			default:
				//当播放者速率跟不上时，会发生丢包
				logs.Debug("camera [%s] http flv sessionId [%d] write fail", hfm.code, wi.GetSessionId())
			}
			return true
		})
	}
}

//添加播放者
func (hfm *HttpFlvManager) AddHttpFlvPlayer(
	playerDone <-chan int,
	pulseInterval time.Duration,
	writer io.Writer,
) (<-chan int, error) {
	sessionId := utils.NextValSnowflakeID()
	//添加缓冲，减少包到达速率震荡导致丢包
	pktStream := make(chan av.Packet, 1024)
	hfw := httpflvwriter.NewHttpFlvWriter(hfm.GetDone(), playerDone, pulseInterval, pktStream, hfm.code, hfm.codecs, writer, sessionId, hfm)
	hfm.hfws.Store(sessionId, hfw)
	go func() {
		<-hfw.GetDone()
		hfm.hfws.Delete(sessionId)
	}()
	return hfw.GetDone(), nil
}

func (hfm *HttpFlvManager) DeleteHFW(sesessionId int64) {
	hfm.hfws.LoadAndDelete(sesessionId)
}
