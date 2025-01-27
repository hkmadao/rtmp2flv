package vo

// 客户端信息
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
