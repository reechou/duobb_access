package models

import (
	"time"
	
	"github.com/reechou/holmes"
)

const (
	PUSH_MSG_MAX = 10
)

type DuobbPushMsg struct {
	ID        int64  `xorm:"pk autoincr"`
	Type      int64  `xorm:"not null default 0 int" json:"type"`
	Msg       string `xorm:"not null default '' varchar(512)" json:"msg"`
	Gray      string `xorm:"not null default '' varchar(256)" json:"gray"`
	CreatedAt int64  `xorm:"not null default 0 int index" json:"createAt"`
}

func CreateDuobbPushMsg(info *DuobbPushMsg) error {
	if info.Type == 0 || info.Msg == "" {
		return CREATE_PUSH_MSG_ERROR_ARGV
	}
	now := time.Now().Unix()
	info.CreatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create duobb push msg error: %v", err)
		return DB_ERROR
	}
	holmes.Info("create duobb push msg[%v] success.", info)
	
	return nil
}

func GetDuobbPushMsgFromTime(t int64) ([]DuobbPushMsg, error) {
	var list []DuobbPushMsg
	err := x.Where("created_at > ?", t).Limit(PUSH_MSG_MAX).Find(&list)
	if err != nil {
		holmes.Error("get push msg list from time error: %v", err)
		return nil, DB_ERROR
	}
	return list, nil
}
