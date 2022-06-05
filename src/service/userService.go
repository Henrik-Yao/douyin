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

func GetUser(userId int64) (model.User, error) {
	var user model.User
	if err := dao.SqlSession.Table("users").Where("id=?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
