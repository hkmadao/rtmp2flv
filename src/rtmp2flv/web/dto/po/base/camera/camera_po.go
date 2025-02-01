package po

import (
	"time"
)

// 摄像头
type CameraPO struct {
	// 摄像头主属性
	Id string `json:"id"`
	// 编号:
	Code string `json:"code"`
	// rtmp识别码:
	RtmpAuthCode string `json:"RtmpAuthCode"`
	// 播放权限码:
	PlayAuthCode string `json:"playAuthCode"`
	// 在线状态:
	OnlineStatus bool `json:"onlineStatus"`
	// 启用状态:
	Enabled bool `json:"enabled"`
	// 保存录像状态:
	SaveVideo bool `json:"saveVideo"`
	// 直播状态:
	Live bool `json:"live"`
	// 创建时间:
	Created time.Time `json:"created"`
	// 加密标志:
	FgSecret bool `json:"fgSecret"`
	// 密钥:
	Secret string `json:"secret"`
	// 被动推送rtmp标志
	FgPassive bool `json:"fgPassive"`
	// 客户端信息主属性:
	IdClientInfo string `json:"idClientInfo"`
	// 客户端信息:
	// ClientInfo ClientInfo `json:"clientInfo"`
	// 摄像头记录
	// CameraRecords []CameraRecord `json:"cameraRecords"`
	// 摄像头分享
	// CameraShares []CameraShare `json:"cameraShares"`
}
