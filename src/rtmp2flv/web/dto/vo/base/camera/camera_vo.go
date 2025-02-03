package vo

import (
	"time"
)

// 摄像头
type CameraVO struct {
	// 摄像头主属性
	Id string `json:"id"`
	// 编号:
	Code string `json:"code"`
	// rtmp识别码:
	RtmpAuthCode string `json:"rtmpAuthCode"`
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
	ClientInfo ClientInfoVO `vo:"ignore" json:"clientInfo"`
}
type ClientInfoVO struct {
	// 客户端信息主属性
	IdClientInfo string `json:"idClientInfo"`
	// 编号:
	ClientCode string `json:"clientCode"`
	// 注册信息签名密钥:
	SignSecret string `json:"signSecret"`
	// 数据传输加密密钥:
	Secret string `json:"secret"`
	// 备注:
	Note string `json:"note"`
}
