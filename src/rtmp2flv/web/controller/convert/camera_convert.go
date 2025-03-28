package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dao/entity"
	camera_po "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dto/po/base/camera"
	camera_vo "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dto/vo/base/camera"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

func ConvertPOToCamera(po camera_po.CameraPO) (camera entity.Camera, err error) {
	err = common.POToEntity(po, &camera)
	if err != nil {
		logs.Error("convertPOToCamera : %v", err)
		err = fmt.Errorf("convertPOToCamera : %v", err)
		return
	}
	return
}

func ConvertPOListToCamera(poes []camera_po.CameraPO) ([]entity.Camera, error) {
	cameras := make([]entity.Camera, len(poes))
	for i, po := range poes {
		camera, err_convert := ConvertPOToCamera(po)
		if err_convert != nil {
			logs.Error("ConvertPOListToCamera : %v", err_convert)
			err := fmt.Errorf("ConvertPOListToCamera : %v", err_convert)
			return nil, err
		}
		cameras[i] = camera
	}
	return cameras, nil
}

func ConvertCameraToVO(camera entity.Camera) (vo camera_vo.CameraVO, err error) {
	vo = camera_vo.CameraVO{}
	err = common.EntityToVO(camera, &vo)
	if err != nil {
		logs.Error("convertCameraToVO : %v", err)
		err = fmt.Errorf("convertCameraToVO : %v", err)
		return
	}
	clientInfo, err := base_service.ClientInfoSelectById(vo.IdClientInfo)
	if err != nil {
		logs.Error("convertCameraToVO : %v", err)
		err = fmt.Errorf("convertCameraToVO : %v", err)
		return
	}
	var clientInfoVO = camera_vo.ClientInfoVO{}
	err = common.EntityToVO(clientInfo, &clientInfoVO)
	if err != nil {
		logs.Error("convertCameraToVO : %v", err)
		err = fmt.Errorf("convertCameraToVO : %v", err)
		return
	}
	vo.ClientInfo = clientInfoVO
	// condition := common.GetEqualCondition("idCamera", vo.Id)
	// var cameraRecordVOList = make([]camera_vo.CameraRecordVO, 0)
	// var cameraRecords = make([]entity.CameraRecord, 0)
	// cameraRecords, err = base_service.CameraRecordFindCollectionByCondition(condition)
	// if err != nil {
	// 	logs.Error("convertCameraToVO : %v", err)
	// 	err = fmt.Errorf("convertCameraToVO : %v", err)
	// 	return
	// }
	// for _, cameraRecord := range cameraRecords {
	// 	var cameraRecordVO = camera_vo.CameraRecordVO{}
	// 	err = common.EntityToVO(cameraRecord, &cameraRecordVO)
	// 	if err != nil {
	// 		logs.Error("convertCameraToVO : %v", err)
	// 		err = fmt.Errorf("convertCameraToVO : %v", err)
	// 		return
	// 	}
	// 	cameraRecordVOList = append(cameraRecordVOList, cameraRecordVO)
	// }
	// vo.cameraRecords = cameraRecordVOList
	// condition := common.GetEqualCondition("cameraId", vo.Id)
	// var cameraShareVOList = make([]camera_vo.CameraShareVO, 0)
	// var cameraShares = make([]entity.CameraShare, 0)
	// cameraShares, err = base_service.CameraShareFindCollectionByCondition(condition)
	// if err != nil {
	// 	logs.Error("convertCameraToVO : %v", err)
	// 	err = fmt.Errorf("convertCameraToVO : %v", err)
	// 	return
	// }
	// for _, cameraShare := range cameraShares {
	// 	var cameraShareVO = camera_vo.CameraShareVO{}
	// 	err = common.EntityToVO(cameraShare, &cameraShareVO)
	// 	if err != nil {
	// 		logs.Error("convertCameraToVO : %v", err)
	// 		err = fmt.Errorf("convertCameraToVO : %v", err)
	// 		return
	// 	}
	// 	cameraShareVOList = append(cameraShareVOList, cameraShareVO)
	// }
	// vo.cameraShares = cameraShareVOList

	return
}

func ConvertCameraToVOList(cameras []entity.Camera) (voList []camera_vo.CameraVO, err error) {
	voList = make([]camera_vo.CameraVO, 0)
	for _, camera := range cameras {
		vo, err_convert := ConvertCameraToVO(camera)
		if err_convert != nil {
			logs.Error("convertCameraToVO : %v", err_convert)
			err = fmt.Errorf("ConvertCameraToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
