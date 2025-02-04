package rtmpserver

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/deepch/vdk/format/rtmp"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dao/entity"
)

// 获取uri信息
func getParamByURI(conn *rtmp.Conn) (string, string, bool) {
	logs.Info("Path : %s , remote port : %s", conn.URL.Path, conn.NetConn().RemoteAddr().String())
	path := conn.URL.Path
	paths := strings.Split(strings.TrimLeft(path, "/"), "/")
	if len(paths) != 2 {
		logs.Error("rtmp path error : %s", path)
		err := conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return "", "", false
	}
	return paths[0], paths[1], true
}

// 权限验证
func authentication(camera entity.Camera, code string, authCode string, conn *rtmp.Conn) bool {
	if camera.RtmpAuthCode != authCode {
		logs.Error("camera %s RtmpAuthCode error : %s", code, authCode)
		conn.Close()
		return false
	}
	if !camera.Enabled {
		logs.Error("camera %s disabled : %s", code, authCode)
		err := conn.Close()
		if err != nil {
			logs.Error("close conn error : %v", err)
		}
		return false
	}
	return true
}

func (r *rtmpServer) Load(key interface{}) (interface{}, bool) {
	return r.rms.Load(key)
}
func (r *rtmpServer) Store(key, value interface{}) {
	r.rms.Store(key, value)
}
func (r *rtmpServer) Delete(key interface{}) {
	r.rms.Delete(key)
}
