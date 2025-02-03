package entity

import (
	"time"
)

// 摄像头
type Camera struct {
	// 摄像头主属性
	Id string `orm:"pk;column(id)" json:"id"`
	// 编号:
	Code string `orm:"column(code)" json:"code"`
	// rtmp识别码:
	RtmpAuthCode string `orm:"column(rtmp_auth_code)" json:"rtmpAuthCode"`
	// 播放权限码:
	PlayAuthCode string `orm:"column(play_auth_code)" json:"playAuthCode"`
	// 在线状态:
	OnlineStatus bool `orm:"column(online_status)" json:"onlineStatus"`
	// 启用状态:
	Enabled bool `orm:"column(enabled)" json:"enabled"`
	// 保存录像状态:
	SaveVideo bool `orm:"column(save_video)" json:"saveVideo"`
	// 直播状态:
	Live bool `orm:"column(live)" json:"live"`
	// 创建时间:
	Created time.Time `orm:"column(created)" json:"created"`
	// 加密标志:
	FgSecret bool `orm:"column(fg_secret)" json:"fgSecret"`
	// 密钥:
	Secret string `orm:"column(secret)" json:"secret"`
	// 被动推送rtmp标志
	FgPassive bool `orm:"column(fg_passive)" json:"fgPassive"`
	// 客户端信息主属性:
	IdClientInfo string `orm:"column(id_client_info)" json:"idClientInfo"`
	// 客户端信息:
	ClientInfo ClientInfo `orm:"-" json:"clientInfo"`
	// 摄像头记录
	CameraRecords []CameraRecord `orm:"-" json:"cameraRecords"`
	// 摄像头分享
	CameraShares []CameraShare `orm:"-" json:"cameraShares"`
}
