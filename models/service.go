package models

import (
	"github.com/reechou/holmes"
)

type Service struct {
	ID          int64  `xorm:"pk autoincr"`
	ServiceName string `xorm:"not null default '' varchar(128)" json:"serviceName"`
	Hosts       string `xorm:"not null default '' varchar(512)" json:"hosts"`
	Type        string `xorm:"not null default '' varchar(64)" json:"type"`
}

type ServiceMethod struct {
	ID            int64  `xorm:"pk autoincr"`
	ServiceMethod string `xorm:"not null default '' varchar(128)" json:"serviceMethod"`
	Uri           string `xorm:"not null default '' varchar(128)" json:"uri"`
}

func CreateService(info *Service) error {
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create duobb access service error: %v", err)
		return DB_ERROR
	}
	holmes.Error("create duobb access service[%v] success.", info)

	return nil
}

func CreateServiceMethod(info *ServiceMethod) error {
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create duobb access service method error: %v", err)
		return DB_ERROR
	}
	holmes.Error("create duobb access service method[%v] success.", info)

	return nil
}

func LoadService() ([]Service, error) {
	var ss []Service
	err := x.Where("id > 0").Find(&ss)
	if err != nil {
		holmes.Error("load service error: %v", err)
		return nil, err
	}
	return ss, nil
}

func LoadServiceMethod() ([]ServiceMethod, error) {
	var sms []ServiceMethod
	err := x.Where("id > 0").Find(&sms)
	if err != nil {
		holmes.Error("load service method error: %v", err)
		return nil, err
	}
	return sms, nil
}
