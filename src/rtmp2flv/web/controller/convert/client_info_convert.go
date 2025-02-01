package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dao/entity"
	client_info_po "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dto/po/base/client_info"
	client_info_vo "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dto/vo/base/client_info"
)

func ConvertPOToClientInfo(po client_info_po.ClientInfoPO) (clientInfo entity.ClientInfo, err error) {
	err = common.POToEntity(po, &clientInfo)
	if err != nil {
		logs.Error("convertPOToClientInfo : %v", err)
		err = fmt.Errorf("convertPOToClientInfo : %v", err)
		return
	}
	return
}

func ConvertPOListToClientInfo(poes []client_info_po.ClientInfoPO) ([]entity.ClientInfo, error) {
	clientInfos := make([]entity.ClientInfo, len(poes))
	for i, po := range poes {
		clientInfo, err_convert := ConvertPOToClientInfo(po)
		if err_convert != nil {
			logs.Error("ConvertPOListToClientInfo : %v", err_convert)
			err := fmt.Errorf("ConvertPOListToClientInfo : %v", err_convert)
			return nil, err
		}
		clientInfos[i] = clientInfo
	}
	return clientInfos, nil
}

func ConvertClientInfoToVO(clientInfo entity.ClientInfo) (vo client_info_vo.ClientInfoVO, err error) {
	vo = client_info_vo.ClientInfoVO{}
	err = common.EntityToVO(clientInfo, &vo)
	if err != nil {
		logs.Error("convertClientInfoToVO : %v", err)
		err = fmt.Errorf("convertClientInfoToVO : %v", err)
		return
	}
	// condition := common.GetEqualCondition("idClientInfo", vo.IdClientInfo)
	// var cameraVOList = make([]client_info_vo.CameraVO, 0)
	// var cameras = make([]entity.Camera, 0)
	// cameras, err = base_service.CameraFindCollectionByCondition(condition)
	// if err != nil {
	// 	logs.Error("convertClientInfoToVO : %v", err)
	// 	err = fmt.Errorf("convertClientInfoToVO : %v", err)
	// 	return
	// }
	// for _, camera := range cameras {
	// 	var cameraVO = client_info_vo.CameraVO{}
	// 	err = common.EntityToVO(camera, &cameraVO)
	// 	if err != nil {
	// 		logs.Error("convertClientInfoToVO : %v", err)
	// 		err = fmt.Errorf("convertClientInfoToVO : %v", err)
	// 		return
	// 	}
	// 	cameraVOList = append(cameraVOList, cameraVO)
	// }
	// vo.cameras = cameraVOList
	
	return
}

func ConvertClientInfoToVOList(clientInfos []entity.ClientInfo) (voList []client_info_vo.ClientInfoVO, err error) {
	voList = make([]client_info_vo.ClientInfoVO, 0)
	for _, clientInfo := range clientInfos {
		vo, err_convert := ConvertClientInfoToVO(clientInfo)
		if err_convert != nil {
			logs.Error("convertClientInfoToVO : %v", err_convert)
			err = fmt.Errorf("ConvertClientInfoToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
