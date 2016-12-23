package models

import (
	"time"

	"github.com/reechou/holmes"
)

type DuobbAccount struct {
	ID        int64  `xorm:"pk autoincr"`
	UserName  string `xorm:"not null default '' varchar(128) index"`
	Password  string `xorm:"not null default '' varchar(128)"`
	Level     int64  `xorm:"not null default 0 int"`
	Status    int64  `xorm:"not null default 0 int"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

func GetDuobbAccount(info *DuobbAccount) error {
	has, err := x.Where("user_name = ?", info.UserName).Get(info)
	if err != nil {
		holmes.Error("get duobb account[%s] error: %v", info.UserName, err)
		return DB_ERROR
	}
	if !has {
		holmes.Error("cannot found account[%s]", info.UserName)
		return GET_ACCOUNT_ERROR_NOTHISMAN
	}
	return nil
}

func CreateDuobbAccount(info *DuobbAccount) error {
	if info.UserName == "" || info.Password == "" {
		return CREATE_ACCOUNT_ERROR_ARGV
	}
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create duobb account error: %v", err)
		return DB_ERROR
	}
	holmes.Info("create duobb account[%v] success.", info)

	return nil
}

func UpdateDuobbAccountPassword(info *DuobbAccount) error {
	if info.UserName == "" || info.Password == "" {
		return UPDATE_ACCOUNT_ERROR_ARGV
	}
	now := time.Now().Unix()
	info.UpdatedAt = now
	_, err := x.Cols("password", "updated_at").Update(info, &DuobbAccount{UserName: info.UserName})
	if err != nil {
		holmes.Error("update duobb account error: %v", err)
		return DB_ERROR
	}

	return nil
}

func UpdateDuobbAccountLevel(info *DuobbAccount) error {
	if info.UserName == "" || info.Password == "" {
		return UPDATE_ACCOUNT_ERROR_ARGV
	}
	now := time.Now().Unix()
	info.UpdatedAt = now
	_, err := x.Cols("level", "updated_at").Update(info, &DuobbAccount{UserName: info.UserName})
	if err != nil {
		holmes.Error("update duobb account error: %v", err)
		return DB_ERROR
	}

	return nil
}
