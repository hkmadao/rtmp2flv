package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/hkmadao/rtmp2flv/src/rtmp2flv/conf" // 必须先导入配置文件
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/rtmpserver"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/task"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web"

	// "net/http"
	// _ "net/http/pprof"

	"github.com/beego/beego/v2/core/logs"
)

func main() {
	rtmpserver.GetSingleRtmpServer().StartRtmpServer()
	task.GetSingleTask().StartTask()
	web.GetSingleWeb().StartWeb()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	logs.Info("Server Start Awaiting Signal")
	// http.ListenAndServe("0.0.0.0:6060", nil)
	select {
	case sig := <-sigs:
		logs.Info(sig)
	case <-done:
	}
	logs.Info("Exiting")
}
