package desc

import (
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
)

func GetClientInfoDesc() *common.EntityDesc {
	var entityInfo = common.EntityInfo{
		Name:        "ClientInfo",
		DisplayName: "客户端信息",
		ClassName:   "ClientInfo",
		TableName:   "client_info",
		BasePath:    "entity::client_info",
	}
	var idClientInfoAttributeInfo = &common.AttributeInfo{
		ColumnName:  "id_client_info",
		Name:        "idClientInfo",
		DisplayName: "客户端信息主属性",
		DataType:    "InternalPK",
		ValueType:   "string",
	}
	var clientCodeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "client_code",
		Name:        "clientCode",
		DisplayName: "编号",
		DataType:    "String",
		ValueType:   "string",
	}
	var signSecretAttributeInfo = &common.AttributeInfo{
		ColumnName:  "sign_secret",
		Name:        "signSecret",
		DisplayName: "注册信息签名密钥",
		DataType:    "String",
		ValueType:   "string",
	}
	var secretAttributeInfo = &common.AttributeInfo{
		ColumnName:  "secret",
		Name:        "secret",
		DisplayName: "数据传输加密密钥",
		DataType:    "String",
		ValueType:   "string",
	}
	var noteAttributeInfo = &common.AttributeInfo{
		ColumnName:  "note",
		Name:        "note",
		DisplayName: "备注",
		DataType:    "String",
		ValueType:   "string",
	}
	var camerasAttributeInfo = &common.AttributeInfo{
		ColumnName:                       "",
		Name:                             "cameras",
		DisplayName:                      "摄像头",
		DataType:                         "InternalArray",
		ValueType:                        "",
		OutEntityName:                    "Camera",
		OutEntityPkAttributeName:         "id",
		OutEntityReversalAttributeName:   "clientInfo",
		OutEntityIdReversalAttributeName: "idClientInfo",
	}
	var entityDesc = &common.EntityDesc{
		EntityInfo:               entityInfo,
		PkAttributeInfo:          idClientInfoAttributeInfo,
		NormalFkIdAttributeInfos: []*common.AttributeInfo{},
		NormalFkAttributeInfos:   []*common.AttributeInfo{},
		NormalChildren: []*common.AttributeInfo{
			camerasAttributeInfo,
		},
		NormalOne2OneChildren: []*common.AttributeInfo{},
		AttributeInfoMap: map[string]*common.AttributeInfo{
			"idClientInfo": idClientInfoAttributeInfo,
			"clientCode":   clientCodeAttributeInfo,
			"signSecret":   signSecretAttributeInfo,
			"secret":       secretAttributeInfo,
			"note":         noteAttributeInfo,
			"cameras":      camerasAttributeInfo,
		},
	}

	return entityDesc
}
