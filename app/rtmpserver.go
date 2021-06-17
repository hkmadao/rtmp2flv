package app

import (
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/hkmadao/rtmp2flv/controllers"
	"github.com/hkmadao/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/services"
)

var rms sync.Map

type RtmpServer struct {
	codeStream <-chan string
}

func NewRtmpServer() *RtmpServer {
	codeStream := controllers.CodeStream()
	rs := &RtmpServer{
		codeStream: codeStream,
	}
	go rs.stopConn()
	go rs.startRtmp()
	return rs
}

func (rs *RtmpServer) stopConn() {

	for {
		code := <-rs.codeStream
		v, b := rms.Load(code)
		if b {
			r := v.(*RtmpManager)
			err := r.conn.Close()
			if err != nil {
				logs.Error("camera [%s] close error : %v", code, err)
				return
			}
			logs.Info("camera [%s] close success", code)
		}
	}

}

func (r *RtmpServer) startRtmp() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	rtmpPort, err := config.Int("server.rtmp.port")
	if err != nil {
		logs.Error("get rtmp port fail : %v", err)
		return
	}
	s := &rtmp.Server{
		Addr:       ":" + strconv.Itoa(rtmpPort),
		HandleConn: handleRtmpConn,
	}
	s.ListenAndServe()
}

func handleRtmpConn(conn *rtmp.Conn) {
	if r := recover(); r != nil {
		logs.Error("HandleConn error : %v", r)
		err := conn.Close()
		if err != nil {
			logs.Error("HandleConn Close err : %v", err)
		}
		return
	}
	NewRtmpManager(conn)
}

type RtmpManager struct {
	conn         *rtmp.Conn
	code         string
	old          bool //是否被挤下线标识
	codecs       []av.CodecData
	done         chan interface{} //告知下游goroutine关闭的chan
	ffmPktStream chan av.Packet
	hfmPktStream chan av.Packet
}

func NewRtmpManager(conn *rtmp.Conn) *RtmpManager {
	r := &RtmpManager{
		conn: conn,
	}
	r.pktTransfer()
	return r
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
	v, b := rms.Load(camera.Code)
	if b {
		logs.Info("camera [%s] online , close old conn", camera.Code)
		oldR := v.(*RtmpManager)
		oldR.old = true
		err = oldR.conn.Close()
		if err != nil {
			logs.Error("camera [%s] close old conn error : %v", camera.Code, err)
		}
	}
	camera.OnlineStatus = 1
	models.CameraUpdate(camera)

	done := make(chan interface{})
	ffmPktChan := make(chan av.Packet, 10)
	hfmPktChan := make(chan av.Packet, 10)

	r.code = camera.Code
	r.codecs = codecs
	r.done = done
	r.ffmPktStream = ffmPktChan
	r.hfmPktStream = hfmPktChan
	r.flvWrite()

	rms.Store(camera.Code, r)

	for {
		pkt, err := r.conn.ReadPacket()
		if err != nil {
			logs.Error("ReadPacket error : %v", err)
			close(done)
			break
		}
		//不能开goroutine,不能保证包的顺序
		r.writeChan(r.done, r.ffmPktStream, pkt)
		r.writeChan(r.done, r.hfmPktStream, pkt)
	}
	//正常掉线
	if !r.old {
		camera, err = models.CameraSelectOne(q)
		if err != nil {
			logs.Error("no camera error : %s", path)
		}
		camera.OnlineStatus = 0
		models.CameraUpdate(camera)

		rms.Delete(r.code)
		err = r.conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
	}
}

func (r *RtmpManager) flvWrite() {
	services.NewHttpFlvManager(r.done, r.hfmPktStream, r.code, r.codecs)

	save, err := config.Bool("server.fileflv.save")
	if err != nil {
		logs.Error("get server.fileflv.save error : %v", err)
		return
	}
	if save {
		services.NewFileFlvManager(r.done, r.ffmPktStream, r.code, r.codecs)
	}
}

func (r *RtmpManager) writeChan(done chan interface{}, pktStream chan<- av.Packet, pkt av.Packet) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	select {
	case pktStream <- pkt:
	case <-time.After(1 * time.Millisecond):
		// logs.Info("lose pkt")
	case <-done:
	}
}
