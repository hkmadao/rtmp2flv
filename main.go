package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/app"
	_ "github.com/hkmadao/rtmp2flv/conf"
)

func main() {
	go app.StartRtmp()
	go app.WebRun()
	go app.ClearToken()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	logs.Info("Server Start Awaiting Signal")
	select {
	case sig := <-sigs:
		logs.Info(sig)
	case <-done:
	}
	logs.Info("Exiting")
}
