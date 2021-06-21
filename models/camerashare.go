package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(CameraShare))
}

type CameraShare struct {
	Id       string    `orm:"pk;column(id)" json:"id"`
	CameraId string    `orm:"column(camera_id)" json:"cameraId"`
	AuthCode string    `orm:"column(auth_code)" json:"authCode"`
	Enabled  int       `orm:"column(enabled)" json:"enabled"`
	Created  time.Time `orm:"column(created)" json:"created"`
}

func CameraShareInsert(e CameraShare) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Insert(&e)
	if err != nil {
		logs.Error("camera insert error : %v", err)
		return i, err
	}
	return i, nil
}

func CameraShareDelete(e CameraShare) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Delete(&e)
	if err != nil {
		logs.Error("camera delete error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraShareUpdate(e CameraShare) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Update(&e)
	if err != nil {
		logs.Error("camera update error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraShareSelectById(id string) (e CameraShare, err error) {
	o := orm.NewOrm()
	e = CameraShare{Id: id}

	err = o.Read(&e)

	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		return e, err
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		return e, err
	} else if err != nil {
		logs.Error("错误: %v", err)
		return e, err
	} else {
		return e, nil
	}
}

func CameraShareSelectOne(q CameraShare) (e CameraShare, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(new(CameraShare)).Filter("cameraId", q.CameraId).Filter("authCode", q.AuthCode).One(&e)
	if err != nil {
		logs.Error("查询出错：%v", err)
		return e, err
	}
	return e, nil
}

func CameraShareCountByCode(code string) (count int64, err error) {
	o := orm.NewOrm()
	count, err = o.QueryTable(new(Camera)).Filter("code", code).Count()
	if err != nil {
		logs.Error("查询出错：%v", err)
		return count, err
	}
	return count, nil
}

func CameraShareSelectAll() (es []CameraShare, err error) {
	o := orm.NewOrm()
	num, err := o.QueryTable(new(Camera)).All(&es)
	if err != nil {
		logs.Error("查询出错：%v", err)
		return es, err
	}
	logs.Info("查询到%d条记录", num)
	return es, nil
}
