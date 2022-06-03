package service

import "errors"

var (
	ErrorUserNameNull   = errors.New("用户名为空")
	ErrorUserNameExtend = errors.New("用户名长度不符合规范")
	ErrorPasswordNull   = errors.New("密码为空")
	ErrorPasswordLength = errors.New("密码长度不符合规范")
	ErrorUserExit       = errors.New("用户已存在")
)
