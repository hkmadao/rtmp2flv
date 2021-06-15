package services

import (
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

type HttpFlvWriter struct {
	sessionId int64
	code      string
	start     bool
	writer    http.ResponseWriter
	codecs    []av.CodecData
	muxer     *flv.Muxer
	close     bool
	done      <-chan interface{}
	heartChan chan<- int
}

func (fw *HttpFlvWriter) IsClose() bool {
	return fw.close
}

func (fw *HttpFlvWriter) HttpWrite(pkt av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("httpWrite panic : %v", r)
		}
	}()
	if fw.start {
		if err := fw.muxer.WritePacket(pkt); err != nil {
			logs.Error("writer packet to httpflv error : %v\n", err)
			close(fw.heartChan)
			fw.close = true
			return
		}
		return
	}
	if pkt.IsKeyFrame {
		muxer := flv.NewMuxer(fw)
		fw.muxer = muxer
		err := fw.muxer.WriteHeader(fw.codecs)
		if err != nil {
			logs.Error("writer header to httpflv error : %v\n", err)
			close(fw.heartChan)
			fw.close = true
			return
		}
		fw.start = true
		if err := fw.muxer.WritePacket(pkt); err != nil {
			logs.Error("writer packet to httpflv error : %v\n", err)
		}
	}

}

//Write extends to io.Writer
func (fw *HttpFlvWriter) Write(p []byte) (n int, err error) {
	n, err = fw.writer.Write(p)
	if err != nil {
		logs.Error("write httpflv error : %v", err)
	}
	for {
		select {
		case fw.heartChan <- 1:
			return
		case <-time.After(1 * time.Millisecond):
			return
		case <-fw.done:
			return
		}
	}

}
