package services

import (
	"os"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/deepch/vdk/av"
)

type FileFlvWriter struct {
	code    string
	isStart bool
	prepare bool
	fd      *os.File
	codecs  []av.CodecData
}

func NewFileFlvWriter(code string, codecs []av.CodecData) *FileFlvWriter {
	return &FileFlvWriter{
		code:    code,
		codecs:  codecs,
		isStart: true,
		prepare: true,
	}
}

func (ffw *FileFlvWriter) Write(p []byte) (n int, err error) {
	if ffw.prepare {
		return
	}
	n, err = ffw.fd.Write(p)
	if err != nil {
		logs.Error("write file error : %v", err)
	}
	return
}
