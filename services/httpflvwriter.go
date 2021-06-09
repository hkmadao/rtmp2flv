package services

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

type HttpFlvWriter struct {
	SessionId int64
	Code      string
	isStart   bool
	Writer    http.ResponseWriter
	codecs    []av.CodecData
	muxer     *flv.Muxer
	close     bool
	errTime   int
	Done      chan<- interface{}
}

//Write extends to io.Writer
func (fw *HttpFlvWriter) Write(p []byte) (n int, err error) {
	n, err = fw.Writer.Write(p)
	if err != nil {
		logs.Error("write httpflv error : %v", err)
	}
	return
}
