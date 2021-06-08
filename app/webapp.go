package app

import (
	"strconv"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gin-gonic/gin"
	"github.com/hkmadao/rtmp2flv/controllers"
)

func WebRun() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("web pain : %v", r)
		}
	}()
	router := gin.Default()
	router.GET("/live/:code/:p.flv", controllers.HttpFlvPlay)

	router.GET("/camera/list", controllers.CameraList)
	router.POST("/camera/edit", controllers.CameraEdit)
	router.POST("/camera/delete/:id", controllers.CameraDelete)

	port, err := config.Int("server.httpflv.port")
	if err != nil {
		logs.Error("get httpflv port error: %v. \n use default port : 8080", err)
		port = 8080
	}
	err = router.Run(":" + strconv.Itoa(port))
	if err != nil {
		logs.Error("Start HTTP Server error", err)
	}
}
