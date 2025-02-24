package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/common"
	"github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dao/entity"
	camera_share_po "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dto/po/base/camera_share"
	camera_share_vo "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/dto/vo/base/camera_share"
	base_service "github.com/hkmadao/rtmp2flv/src/rtmp2flv/web/service/base"
)

func ConvertPOToCameraShare(po camera_share_po.CameraSharePO) (cameraShare entity.CameraShare, err error) {
	err = common.POToEntity(po, &cameraShare)
	if err != nil {
		logs.Error("convertPOToCameraShare : %v", err)
		err = fmt.Errorf("convertPOToCameraShare : %v", err)
		return
	}
	return
}

func ConvertPOListToCameraShare(poes []camera_share_po.CameraSharePO) ([]entity.CameraShare, error) {
	cameraShares := make([]entity.CameraShare, len(poes))
	for i, po := range poes {
		cameraShare, err_convert := ConvertPOToCameraShare(po)
		if err_convert != nil {
			logs.Error("ConvertPOListToCameraShare : %v", err_convert)
			err := fmt.Errorf("ConvertPOListToCameraShare : %v", err_convert)
			return nil, err
		}
		cameraShares[i] = cameraShare
	}
	return cameraShares, nil
}

func ConvertCameraShareToVO(cameraShare entity.CameraShare) (vo camera_share_vo.CameraShareVO, err error) {
	vo = camera_share_vo.CameraShareVO{}
	err = common.EntityToVO(cameraShare, &vo)
	if err != nil {
		logs.Error("convertCameraShareToVO : %v", err)
		err = fmt.Errorf("convertCameraShareToVO : %v", err)
		return
	}
	camera, err := base_service.CameraSelectById(vo.CameraId)
	if err != nil {
		logs.Error("convertCameraShareToVO : %v", err)
		err = fmt.Errorf("convertCameraShareToVO : %v", err)
		return
	}
	var cameraVO = camera_share_vo.CameraVO{}
	err = common.EntityToVO(camera, &cameraVO)
	if err != nil {
		logs.Error("convertCameraShareToVO : %v", err)
		err = fmt.Errorf("convertCameraShareToVO : %v", err)
		return
	}
	vo.Camera = cameraVO

	return
}

func ConvertCameraShareToVOList(cameraShares []entity.CameraShare) (voList []camera_share_vo.CameraShareVO, err error) {
	voList = make([]camera_share_vo.CameraShareVO, 0)
	for _, cameraShare := range cameraShares {
		vo, err_convert := ConvertCameraShareToVO(cameraShare)
		if err_convert != nil {
			logs.Error("convertCameraShareToVO : %v", err_convert)
			err = fmt.Errorf("ConvertCameraShareToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
