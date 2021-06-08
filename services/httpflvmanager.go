package services

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

var Hms map[string]*HttpFlvManager

func init() {
	Hms = make(map[string]*HttpFlvManager)
}

type HttpFlvManager struct {
	codecs []av.CodecData
	Fws    map[string]*HttpFlvWriter
}

func NewHttpFlvManager() *HttpFlvManager {
	hm := &HttpFlvManager{
		Fws: make(map[string]*HttpFlvWriter),
	}
	return hm
}

func (fm *HttpFlvManager) codec(code string, codecs []av.CodecData) {
	fm.codecs = codecs
	Hms[code] = fm
}

//Write extends to writer.Writer
func (fm *HttpFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("HttpFlvManager FlvWrite pain %v", r)
		}
	}()
	fm.codec(code, codecs)
	for {
		select {
		case <-done:
			return
		case pkt := <-pchan:
			for _, fw := range fm.Fws {
				if fw.close {
					fw.Done <- nil
					delete(fm.Fws, fw.SessionId)
					continue
				}
				if fw.isStart {
					if err := fw.muxer.WritePacket(pkt); err != nil {
						logs.Error("writer packet to httpflv error : %v\n", err)
						if fw.errTime > 20 {
							fw.close = true
							continue
						}
						fw.errTime = fw.errTime + 1
					} else {
						fw.errTime = 0
					}
					continue
				}
				if pkt.IsKeyFrame {
					muxer := flv.NewMuxer(fw)
					fw.muxer = muxer
					err := fw.muxer.WriteHeader(fm.codecs)
					if err != nil {
						logs.Error("writer header to httpflv error : %v\n", err)
						if fw.errTime > 20 {
							fw.close = true
							continue
						}
						fw.errTime = fw.errTime + 1
					}
					fw.isStart = true
					if err := fw.muxer.WritePacket(pkt); err != nil {
						logs.Error("writer packet to httpflv error : %v\n", err)
					}
				}
			}
		}
	}
}
