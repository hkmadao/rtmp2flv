package desc

import (
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
)

func GetCameraDesc() *common.EntityDesc {
	var entityInfo = common.EntityInfo{
		Name:        "Camera",
		DisplayName: "摄像头",
		ClassName:   "Camera",
		TableName:   "camera",
		BasePath:    "entity::camera",
	}
	var idAttributeInfo = &common.AttributeInfo{
		ColumnName:  "id",
		Name:        "id",
		DisplayName: "摄像头主属性",
		DataType:    "InternalPK",
		ValueType:   "string",
	}
	var codeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "code",
		Name:        "code",
		DisplayName: "编号",
		DataType:    "String",
		ValueType:   "string",
	}
	var RtmpAuthCodeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "rtmp_auth_code",
		Name:        "RtmpAuthCode",
		DisplayName: "rtmp识别码",
		DataType:    "String",
		ValueType:   "string",
	}
	var playAuthCodeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "play_auth_code",
		Name:        "playAuthCode",
		DisplayName: "播放权限码",
		DataType:    "String",
		ValueType:   "string",
	}
	var onlineStatusAttributeInfo = &common.AttributeInfo{
		ColumnName:  "online_status",
		Name:        "onlineStatus",
		DisplayName: "在线状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var enabledAttributeInfo = &common.AttributeInfo{
		ColumnName:  "enabled",
		Name:        "enabled",
		DisplayName: "启用状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var saveVideoAttributeInfo = &common.AttributeInfo{
		ColumnName:  "save_video",
		Name:        "saveVideo",
		DisplayName: "保存录像状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var liveAttributeInfo = &common.AttributeInfo{
		ColumnName:  "live",
		Name:        "live",
		DisplayName: "直播状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var createdAttributeInfo = &common.AttributeInfo{
		ColumnName:  "created",
		Name:        "created",
		DisplayName: "创建时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var fgSecretAttributeInfo = &common.AttributeInfo{
		ColumnName:  "fg_secret",
		Name:        "fgSecret",
		DisplayName: "加密标志",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var secretAttributeInfo = &common.AttributeInfo{
		ColumnName:  "secret",
		Name:        "secret",
		DisplayName: "密钥",
		DataType:    "String",
		ValueType:   "string",
	}
	var fgPassiveAttributeInfo = &common.AttributeInfo{
		ColumnName:  "fg_passive",
		Name:        "fgPassive",
		DisplayName: "被动推送rtmp标志",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var idClientInfoAttributeInfo = &common.AttributeInfo{
		ColumnName:                     "id_client_info",
		Name:                           "idClientInfo",
		DisplayName:                    "客户端信息主属性",
		DataType:                       "InternalFK",
		ValueType:                      "string",
		InnerAttributeName:             "clientInfo",
		OutEntityName:                  "ClientInfo",
		OutEntityPkAttributeName:       "idClientInfo",
		OutEntityReversalAttributeName: "cameras",
	}
	var clientInfoAttributeInfo = &common.AttributeInfo{
		ColumnName:                     "",
		Name:                           "clientInfo",
		DisplayName:                    "客户端信息",
		DataType:                       "InternalRef",
		ValueType:                      "",
		InnerAttributeName:             "idClientInfo",
		OutEntityName:                  "ClientInfo",
		OutEntityPkAttributeName:       "idClientInfo",
		OutEntityReversalAttributeName: "cameras",
	}
	var cameraRecordsAttributeInfo = &common.AttributeInfo{
		ColumnName:                       "",
		Name:                             "cameraRecords",
		DisplayName:                      "摄像头记录",
		DataType:                         "InternalArray",
		ValueType:                        "",
		OutEntityName:                    "CameraRecord",
		OutEntityPkAttributeName:         "idCameraRecord",
		OutEntityReversalAttributeName:   "camera",
		OutEntityIdReversalAttributeName: "idCamera",
	}
	var cameraSharesAttributeInfo = &common.AttributeInfo{
		ColumnName:                       "",
		Name:                             "cameraShares",
		DisplayName:                      "摄像头分享",
		DataType:                         "InternalArray",
		ValueType:                        "",
		OutEntityName:                    "CameraShare",
		OutEntityPkAttributeName:         "id",
		OutEntityReversalAttributeName:   "camera",
		OutEntityIdReversalAttributeName: "cameraId",
	}
	var entityDesc = &common.EntityDesc{
		EntityInfo:      entityInfo,
		PkAttributeInfo: idAttributeInfo,
		NormalFkIdAttributeInfos: []*common.AttributeInfo{
			idClientInfoAttributeInfo,
		},
		NormalFkAttributeInfos: []*common.AttributeInfo{
			clientInfoAttributeInfo,
		},
		NormalChildren: []*common.AttributeInfo{
			cameraRecordsAttributeInfo,
			cameraSharesAttributeInfo,
		},
		NormalOne2OneChildren: []*common.AttributeInfo{},
		AttributeInfoMap: map[string]*common.AttributeInfo{
			"id":            idAttributeInfo,
			"code":          codeAttributeInfo,
			"RtmpAuthCode":  RtmpAuthCodeAttributeInfo,
			"playAuthCode":  playAuthCodeAttributeInfo,
			"onlineStatus":  onlineStatusAttributeInfo,
			"enabled":       enabledAttributeInfo,
			"saveVideo":     saveVideoAttributeInfo,
			"live":          liveAttributeInfo,
			"created":       createdAttributeInfo,
			"fgSecret":      fgSecretAttributeInfo,
			"secret":        secretAttributeInfo,
			"fgPassive":     fgPassiveAttributeInfo,
			"idClientInfo":  idClientInfoAttributeInfo,
			"clientInfo":    clientInfoAttributeInfo,
			"cameraRecords": cameraRecordsAttributeInfo,
			"cameraShares":  cameraSharesAttributeInfo,
		},
	}

	return entityDesc
}
