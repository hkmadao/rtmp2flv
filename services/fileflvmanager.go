package services

import (
	"os"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/core/config"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/flv"
)

type FileFlvManager struct {
	fw *FileFlvWriter
}

func NewFileFlvManager() *FileFlvManager {
	return &FileFlvManager{}
}

func (fm *FileFlvManager) codec(code string, codecs []av.CodecData) {
	fd, err := os.OpenFile(getFileFlvPath()+"/"+code+"_"+time.Now().Format("20060102150405")+".flv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logs.Error("open file error :", err)
	}
	fm.fw = &FileFlvWriter{
		codecs: codecs,
		code:   code,
		fd:     fd,
	}
}

//Write extends to writer.Writer
func (fm *FileFlvManager) FlvWrite(code string, codecs []av.CodecData, done <-chan interface{}, pchan <-chan av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("FlvFileManager FlvWrite panic %v", r)
		}
	}()
	fm.codec(code, codecs)
	muxer := flv.NewMuxer(fm.fw)
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-done:
			fm.fw.fd.Close()
			return
		case <-ticker.C: //split flvFile
			fd, err := os.OpenFile(getFileFlvPath()+"/"+fm.fw.code+"_"+time.Now().Format("20060102150405")+".flv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				logs.Error("open file error :", err)
			}
			fdOld := fm.fw.fd
			fm.fw.prepare = true
			fm.fw.isStart = false
			fm.fw.fd = fd
			fm.fw.prepare = false
			fdOld.Close()
		case pkt := <-pchan:
			if fm.fw.isStart {
				if err := muxer.WritePacket(pkt); err != nil {
					logs.Error("writer packet to flv file error : %v\n", err)
				}
				continue
			}
			if pkt.IsKeyFrame {
				err := muxer.WriteHeader(fm.fw.codecs)
				if err != nil {
					logs.Error("writer header to flv file error : %v\n", err)
				}
				if err := muxer.WritePacket(pkt); err != nil {
					logs.Error("writer packet to flv file error : %v\n", err)
				}
				fm.fw.isStart = true
			}
		}
	}
}

func getFileFlvPath() string {
	fileFlvPath, err := config.String("server.fileflv.path")
	if err != nil {
		logs.Error("get fileflv path error :", err)
		return ""
	}
	return fileFlvPath
}
