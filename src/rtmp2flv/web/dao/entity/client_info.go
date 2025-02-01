package entity

// 客户端信息
type ClientInfo struct {
	// 客户端信息主属性
	IdClientInfo string `orm:"pk;column(id_client_info)" json:"idClientInfo"`
	// 编号:
	ClientCode string `orm:"column(client_code)" json:"clientCode"`
	// 注册信息签名密钥:
	SignSecret string `orm:"column(sign_secret)" json:"signSecret"`
	// 数据传输加密密钥:
	Secret string `orm:"column(secret)" json:"secret"`
	// 备注:
	Note string `orm:"column(note)" json:"note"`
	// 摄像头
	Cameras []Camera `orm:"-" json:"cameras"`
}
