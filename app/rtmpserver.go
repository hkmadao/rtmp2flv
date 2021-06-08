package app

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/hkmadao/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/services"
)

func StartRtmp() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("RTMP server panic: ", r)
		}
	}()
	rtmpPort, err := config.Int("server.rtmp.port")
	if err != nil {
		logs.Error("get rtmp port fail : %v", err)
		return
	}
	s := &rtmp.Server{
		Addr:       ":" + strconv.Itoa(rtmpPort),
		HandleConn: HandleConn,
	}
	s.ListenAndServe()
}

func HandleConn(conn *rtmp.Conn) {
	if r := recover(); r != nil {
		logs.Error("HandleConn error : %v", r)
		err := conn.Close()
		if err != nil {
			logs.Error("HandleConn Close err : %v", err)
		}
		return
	}
	err := conn.Prepare()
	if err != nil {
		logs.Error("Prepare error : %v", err)
		return
	}
	//权限验证
	logs.Info("Path : %s", conn.URL.Path)
	path := conn.URL.Path
	paths := strings.Split(strings.TrimLeft(path, "/"), "/")
	if len(paths) != 2 {
		logs.Error("rtmp path error : %s", path)
		conn.Close()
		return
	}
	q := models.Camera{Code: paths[0]}
	camera, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("no camera error : %s", path)
		conn.Close()
		return
	}
	if camera.RtmpAuthCode != paths[1] {
		logs.Error("RtmpAuthCode error : %s", path)
		conn.Close()
		return
	}
	if camera.Enabled != 1 {
		logs.Error("camera disabled : %s", path)
		conn.Close()
		return
	}

	codecs, err := conn.Streams()
	if err != nil {
		logs.Error("get codecs error : %v", err)
		conn.Close()
		return
	}

	done := make(chan interface{})
	ffmPktChan := make(chan av.Packet)
	hfmPktChan := make(chan av.Packet)
	ffm := services.NewFileFlvManager()
	hfm := services.NewHttpFlvManager()
	go ffm.FlvWrite(camera.Code, codecs, done, ffmPktChan)
	go hfm.FlvWrite(camera.Code, codecs, done, hfmPktChan)

	for {
		pkt, err := conn.ReadPacket()
		if err != nil {
			logs.Error("ReadPacket error : %v", err)
			close(done)
			break
		}
		select {
		case ffmPktChan <- pkt:
		case <-time.After(1 * time.Microsecond):
		}
		select {
		case hfmPktChan <- pkt:
		case <-time.After(1 * time.Microsecond):
		}
	}
	conn.Close()
}
