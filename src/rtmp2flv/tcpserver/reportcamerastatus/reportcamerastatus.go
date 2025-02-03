package reportcamerastatus

import (
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

var lastReportInfo sync.Map

func OnlineStatus(cameraCode string) {
	condition := common.GetEqualCondition("code", cameraCode)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("code: %s find camere error: %v", cameraCode)
		return
	}
	if !camera.Enabled {
		return
	}
	camera.OnlineStatus = true
	lastReportInfo.Store(cameraCode, time.Now())
	_, err = base_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("code: %s update camere onlieStatus error: %v", cameraCode)
	}
}

func OfflineStatus(cameraCode string) {
	condition := common.GetEqualCondition("code", cameraCode)
	camera, err := base_service.CameraFindOneByCondition(condition)
	if err != nil {
		logs.Error("code: %s find camere error: %v", cameraCode)
		return
	}
	if !camera.Enabled {
		return
	}
	camera.OnlineStatus = false
	_, err = base_service.CameraUpdateById(camera)
	if err != nil {
		logs.Error("code: %s update camere onlieStatus error: %v", cameraCode)
	}
}

func CheckExpires(cameraCode string) bool {
	value, ok := lastReportInfo.Load(cameraCode)
	if !ok {
		return false
	}
	lastReportTime := value.(time.Time)
	return time.Since(lastReportTime) < 3*time.Minute
}
