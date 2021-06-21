package app

import (
	"encoding/json"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/models"
	"github.com/hkmadao/rtmp2flv/result"
	"github.com/hkmadao/rtmp2flv/services"
)

func HttpFlvServe() error {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	port, err := config.Int("server.httpflv.port")
	if err != nil {
		logs.Error("get httpflv port error: %v. \n use default port : 9091", err)
		port = 9091
	}
	httpflvAddr := ":" + strconv.Itoa(port)
	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		logs.Error("%v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/live/", handleConn)
	if err := http.Serve(flvListen, mux); err != nil {
		return err
	}
	return nil
}

func handleConn(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")
	uri := strings.TrimSuffix(strings.TrimLeft(req.RequestURI, "/"), ".flv")
	uris := strings.Split(uri, "/")
	if len(uris) < 3 || uris[0] != "live" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	method := uris[1]
	code := uris[2]
	authCode := uris[3]
	r := result.Result{
		Code: 1,
		Msg:  "",
	}
	q := models.Camera{Code: code}
	camera, err := models.CameraSelectOne(q)
	if err != nil {
		logs.Error("camera query error : %v", err)
		r.Code = 0
		r.Msg = "camera query error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if !(method == "temp" || method == "permanent") {
		logs.Error("method error : %s", method)
		r.Code = 0
		r.Msg = "method error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if method == "temp" {
		csq := models.CameraShare{CameraId: camera.Id, AuthCode: authCode}
		cs, err := models.CameraShareSelectOne(csq)
		if err != nil {
			logs.Error("CameraShareSelectOne error : %v", err)
			r.Code = 0
			r.Msg = "system error"
			rbytes, _ := json.Marshal(r)
			w.Write(rbytes)
			return
		}
		if time.Now().After(cs.Created.Add(7 * 24 * time.Hour)) {
			logs.Error("camera [%s] AuthCodeTemp expired : %s", camera.Code, authCode)
			r.Code = 0
			r.Msg = "authCode expired"
			rbytes, _ := json.Marshal(r)
			w.Write(rbytes)
			return
		}

	}
	if method == "permanent" && authCode != camera.PlayAuthCode {
		logs.Error("AuthCodePermanent error : %s", authCode)
		r.Code = 0
		r.Msg = "authCode error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	if !services.ExistsHttpFlvManager(code) {
		logs.Error("camera [%s] no connection : %s", code)
		r.Code = 0
		r.Msg = "camera no connection"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
	logs.Info("player [%s] addr [%s] connecting", code, req.RemoteAddr)
	//管理员可以主动中断播放
	endStream, heartbeatStream, _, err := services.AddHttpFlvPlayer(code, w)
	if err != nil {
		logs.Error("camera [%s] add player error : %s", code)
		r.Code = 0
		r.Msg = "add player error"
		rbytes, _ := json.Marshal(r)
		w.Write(rbytes)
		return
	}
Loop:
	for {
		select {
		case <-endStream:
			break Loop
		case <-heartbeatStream:
			continue
		case <-time.After(10 * time.Second):
			logs.Info("player [%s] addr [%s] timeout exit", code, req.RemoteAddr)
			break Loop
		}
	}
	logs.Info("player [%s] addr [%s] exit", code, req.RemoteAddr)
}
