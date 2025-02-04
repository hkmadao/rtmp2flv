package task

import (
	"runtime/debug"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvadmin"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/flvadmin/fileflvmanager/fileflvreader"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/rtmpserver"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/tcpserver/reportcamerastatus"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

var taskInstance *task

func init() {
	taskInstance = &task{}
}

type task struct {
}

func GetSingleTask() *task {
	return taskInstance
}

func (t *task) StartTask() {
	go t.clearToken()
	go t.offlineCamera()
	go t.ClearHistoryVideo()
	go t.SendStopRtmpPush()
}

func (t *task) clearToken() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		<-time.After(24 * time.Hour)
	}
}

func (t *task) offlineCamera() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		condition := common.GetEqualCondition("onlineStatus", true)
		css, err := base_service.CameraFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("query camera error : %v", err)
		}
		for _, cs := range css {
			if cs.FgPassive {
				if expires := reportcamerastatus.CheckExpires(cs.Code); !expires {
					cs.OnlineStatus = false
					base_service.CameraUpdateById(cs)
				}
				continue
			}
			if cs.FgSecret {
				if exists := rtmpserver.GetSingleEncryptRtmpServer().ExistsPublisher(cs.Code); !exists {
					cs.OnlineStatus = false
					base_service.CameraUpdateById(cs)
				}
				continue
			}
			if exists := rtmpserver.GetSingleRtmpServer().ExistsPublisher(cs.Code); !exists {
				cs.OnlineStatus = false
				base_service.CameraUpdateById(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}

func (t *task) ClearHistoryVideo() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		ltCreated := time.Now().Add(-7 * 24 * time.Hour)
		condition := common.GetLtCondition("created", ltCreated.Format(time.RFC3339))
		css, err := base_service.CameraRecordFindCollectionByCondition(condition)
		if err != nil {
			logs.Error("query CameraRecord error : %v", err)
		}
		for _, cs := range css {
			fileExists := false
			if cs.FgTemp {
				fileExists = fileflvreader.FlvFileExists(cs.TempFileName)
			} else {
				fileExists = fileflvreader.FlvFileExists(cs.FileName)
			}
			if !fileExists {
				logs.Info("flv file: %s not exists, clean", cs.FileName)
				base_service.CameraRecordDelete(cs)
			}
		}
		<-time.After(10 * time.Minute)
	}
}

func (t *task) SendStopRtmpPush() {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("system painc : %v \nstack : %v", r, string(debug.Stack()))
		}
	}()
	for {
		flvadmin.GetSingleHttpFlvAdmin().TickerCheckStopRtmp()
		<-time.After(10 * time.Minute)
	}
}
