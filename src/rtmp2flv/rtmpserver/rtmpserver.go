package rtmpserver

import (
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/rtmpserver/rtmppublisher"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	ext_controller "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/controller/ext"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

var rtmpserverInstance *rtmpServer

func init() {
	rtmpserverInstance = &rtmpServer{}
}

type rtmpServer struct {
	rms   sync.Map
	conns sync.Map
}

func GetSingleRtmpServer() *rtmpServer {
	return rtmpserverInstance
}

func (rs *rtmpServer) StartRtmpServer() {
	go rs.startRtmp()
	done := make(chan interface{})
	go rs.stopConn(done, ext_controller.CodeStream())
}

func (rs *rtmpServer) ExistsPublisher(code string) bool {
	exists := false
	rs.rms.Range(func(key, value interface{}) bool {
		codeKey := key.(string)
		if code == codeKey {
			exists = true
			return false
		}
		return true
	})
	return exists
}

func (rs *rtmpServer) stopConn(done <-chan interface{}, codeStream <-chan string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		select {
		case <-done:
			return
		case code := <-codeStream:
			v, b := rs.conns.Load(code)
			if b {
				r := v.(*rtmp.Conn)
				err := r.Close()
				if err != nil {
					logs.Error("camera [%s] close error : %v", code, err)
					continue
				}
				logs.Info("camera [%s] close success", code)
			}
		}
	}

}

func (r *rtmpServer) startRtmp() {
	defer func() {
		if recover_rusult := recover(); recover_rusult != nil {
			logs.Error("system painc : %v \nstack : %v", recover_rusult, string(debug.Stack()))
		}
	}()
	rtmpPort, err := config.Int("server.rtmp.port")
	if err != nil {
		logs.Error("get rtmp port fail : %v", err)
		return
	}
	// rtmp.Debug = true
	s := &rtmp.Server{
		Addr:       ":" + strconv.Itoa(rtmpPort),
		HandleConn: r.handleRtmpConn,
	}
	if err := s.ListenAndServe(); err != nil {
		logs.Error("encrypt rtmp ListenAndServe error: %v", err)
	}
}

func (r *rtmpServer) handleRtmpConn(conn *rtmp.Conn) {
	defer func() {
		if recover_rusult := recover(); recover_rusult != nil {
			logs.Error("HandleConn error : %v", recover_rusult)
		}
	}()
	defer func() {
		err := conn.Close()
		if err != nil {
			logs.Error("HandleConn Close err : %v", err)
		}
	}()
	logs.Info("client arrive : %s", conn.NetConn().RemoteAddr().String())
	err := conn.Prepare()
	if err != nil {
		logs.Error("Prepare error : %v , remote port : %s", err, conn.NetConn().RemoteAddr().String())
		err = conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}

	code, authCode, ok := getParamByURI(conn)
	if !ok {
		return
	}

	condition := common.GetEqualCondition("code", code)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("find camera by code: %s error : %v", code, err)
		return
	}

	if camera.FgEncrypt {
		logs.Error("camera: %s fgEncrypt is %b", code, camera.FgEncrypt)
		return
	}

	if ok := authentication(camera, code, authCode, conn); !ok {
		return
	}

	logs.Info("publish authentication success : %s", code)

	codecs, err := conn.Streams()
	if err != nil {
		logs.Error("get codecs error : %v", err)
		err = conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return
	}
	v, ok := r.conns.LoadAndDelete(camera.Code)
	if ok {
		logs.Info("camera [%s] online , close old conn", code)
		conn := v.(*rtmp.Conn)
		err := conn.Close()
		if err != nil {
			logs.Error("camera [%s] close old conn error : %v", code, err)
		}
	}
	v, ok = r.rms.Load(camera.Code)
	if ok {
		logs.Info("camera [%s] online , close old conn", camera.Code)
		oldR := v.(*rtmppublisher.Publisher)
		//等待旧连接关闭完成
		oldR.Done()
	}
	r.conns.Store(camera.Code, conn)

	camera.OnlineStatus = true
	base_service.CameraUpdateById(camera)

	done := make(chan interface{})
	//添加缓冲
	pktStream := make(chan av.Packet, 1024)
	defer func() {
		close(done)
		close(pktStream)
	}()

	p := rtmppublisher.NewPublisher(done, pktStream, code, codecs, r)
	r.rms.Store(camera.Code, p)
	for {
		pkt, err := conn.ReadPacket()
		if err != nil {
			logs.Error("ReadPacket error : %v", err)
			break
		}
		select {
		case pktStream <- pkt:
		default:
			//添加缓冲，缓解前后速率不一致问题，但是如果收包平均速率大于消费平均速率，依然会导致丢包
			logs.Debug("rtmpserver lose packet")
		}
	}

	camera, err = base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("no camera error : %s", code)
	} else {
		if !camera.FgPassive {
			camera.OnlineStatus = false
			base_service.CameraUpdateById(camera)
		}
	}

	r.rms.Delete(code)
	r.conns.Delete(code)
	err = conn.Close()
	if err != nil {
		logs.Error("close conn error : %v", err)
	}

}
