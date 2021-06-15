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
	r := &RtmpManager{
		conn: conn,
	}
	r.pktTransfer()
}

type RtmpManager struct {
	conn       *rtmp.Conn
	code       string
	codecs     []av.CodecData
	done       chan interface{}
	ffmPktChan chan av.Packet
	hfmPktChan chan av.Packet
}

func (r *RtmpManager) pktTransfer() {
	err := r.conn.Prepare()
	if err != nil {
		logs.Error("Prepare error : %v , remote port : %s", err, r.conn.NetConn().RemoteAddr().String())
		err = r.conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}
	//权限验证
	logs.Info("Path : %s , remote port : %s", r.conn.URL.Path, r.conn.NetConn().RemoteAddr().String())
	path := r.conn.URL.Path
	paths := strings.Split(strings.TrimLeft(path, "/"), "/")
	if len(paths) != 2 {
		logs.Error("rtmp path error : %s", path)
		err = r.conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}
	q := models.Camera{Code: paths[0]}
	camera, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("no camera error : %s", path)
		err = r.conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}
	if camera.RtmpAuthCode != paths[1] {
		logs.Error("RtmpAuthCode error : %s", path)
		r.conn.Close()
		return
	}
	if camera.Enabled != 1 {
		logs.Error("camera disabled : %s", path)
		err = r.conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}

	codecs, err := r.conn.Streams()
	if err != nil {
		logs.Error("get codecs error : %v", err)
		err = r.conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}
	camera.OnlineStatus = 1
	models.CameraUpdate(camera)

	done := make(chan interface{})
	ffmPktChan := make(chan av.Packet, 10)
	hfmPktChan := make(chan av.Packet, 10)

	r.code = camera.Code
	r.codecs = codecs
	r.done = done
	r.ffmPktChan = ffmPktChan
	r.hfmPktChan = hfmPktChan
	r.flvWrite()

	for {
		pkt, err := r.conn.ReadPacket()
		if err != nil {
			logs.Error("ReadPacket error : %v", err)
			close(done)
			break
		}
		r.writeChan(pkt)
	}
	err = r.conn.Close()
	if err != nil {
		logs.Error("close conn error : %v", err)
	}
}

func (r *RtmpManager) flvWrite() {
	hfm := services.NewHttpFlvManager()
	go hfm.FlvWrite(r.code, r.codecs, r.done, r.hfmPktChan)

	save, err := config.Bool("server.fileflv.save")
	if err != nil {
		logs.Error("get server.fileflv.save error : %v", err)
		return
	}
	if save {
		ffm := services.NewFileFlvManager()
		go ffm.FlvWrite(r.code, r.codecs, r.done, r.ffmPktChan)
	}
}

func (r *RtmpManager) writeChan(pkt av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("writeChan panicc : %v", r)
		}
	}()
	select {
	case r.ffmPktChan <- pkt:
	case <-time.After(1 * time.Millisecond):
	case <-r.done:
	}
	select {
	case r.hfmPktChan <- pkt:
	case <-time.After(1 * time.Millisecond):
	case <-r.done:
	}
}
