package service

import (
	"douyin/src/dao"
	"douyin/src/model"
)

func CreateUser(user *model.User) (err error) {
	if err = dao.SqlSession.Create(user).Error; err != nil {
		return err
	}
	return
}
