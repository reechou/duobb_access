package models

import (
	"errors"
)

var (
	DB_ERROR                    = errors.New("数据库错误")
	CREATE_ACCOUNT_ERROR_ARGV   = errors.New("创建账户参数错误")
	UPDATE_ACCOUNT_ERROR_ARGV   = errors.New("创建账户参数错误")
	GET_ACCOUNT_ERROR_NOTHISMAN = errors.New("无此账号")
	CREATE_PUSH_MSG_ERROR_ARGV  = errors.New("创建账户参数错误")
)
