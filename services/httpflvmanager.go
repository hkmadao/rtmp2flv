package services

import (
	"net/http"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/deepch/vdk/av"
	"github.com/hkmadao/rtmp2flv/utils"
)

var hfms map[string]*HttpFlvManager

func init() {
	hfms = make(map[string]*HttpFlvManager)
}

//添加播放者
func AddHttpFlvPlayer(done <-chan interface{}, code string, writer http.ResponseWriter) <-chan int {
	heartChan := make(chan int)
	sessionId := utils.NextValSnowflakeID()
	fw := &HttpFlvWriter{
		sessionId: sessionId,
		writer:    writer,
		heartChan: heartChan,
		codecs:    hfms[code].codecs,
		code:      code,
	}
	hfms[code].fms[sessionId] = fw
	return heartChan
}

type HttpFlvManager struct {
	codecs []av.CodecData
	fms    map[int64]*HttpFlvWriter
}

func NewHttpFlvManager() *HttpFlvManager {
	hm := &HttpFlvManager{}
	return hm
}

func (hfm *HttpFlvManager) codec(code string, codecs []av.CodecData) {
	hfm.fms = make(map[int64]*HttpFlvWriter)
	hfm.codecs = codecs
	hfms[code] = hfm
}

//Write extends to writer.Writer
func (hfm *HttpFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("HttpFlvManager FlvWrite panic %v", r)
		}
	}()
	hfm.codec(code, codecs)
	for {
		select {
		case <-done:
			return
		case pkt := <-pchan:
			deleteKeys := make([]int64, 2)
			for _, fw := range hfm.fms {
				if fw.IsClose() {
					deleteKeys = append(deleteKeys, fw.sessionId)
				}
				go fw.HttpWrite(pkt)
			}
			for _, sessionId := range deleteKeys {
				delete(hfm.fms, sessionId)
			}
		}
	}
}
