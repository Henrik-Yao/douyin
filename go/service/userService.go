package service

import (
	"douyin/go/dao"
	"douyin/go/model"
)

func CreateUser(user *model.User) (err error) {
	if err = dao.SqlSession.Create(user).Error; err != nil {
		return err
	}
	return
}
