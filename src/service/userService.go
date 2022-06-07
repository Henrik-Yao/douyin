package service

import (
	"crypto/md5"
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	MaxUsernameLength = 32            //用户名最大长度
	MaxPasswordLength = 32            //密码最大长度
	MinPasswordLength = 6             //密码最小长度
	secret            = "return11111" //密码加密
)

//增

func CreateRegisterUser(userName string, passWord string) (model.User, error) {
	//1.Following数据模型准备
	// 记录一下原始密码(用户登录的密码)
	originPassword := passWord
	// 新生成加密密码用于和查询到的密码比较
	newPassword := EncryptPassword([]byte(originPassword))
	newUser := model.User{
		Name:     userName,
		Password: newPassword,
	}
	//2.模型关联到数据库表users //可注释
	dao.SqlSession.AutoMigrate(&model.User{})
	//3.新建user
	if IsUserExistByName(userName) {
		//用户已存在
		return newUser, common.ErrorUserExit
	} else {
		//用户不存在，新建用户
		if err := dao.SqlSession.Model(&model.User{}).Create(&newUser).Error; err != nil {
			//错误处理
			fmt.Println(err)
			return newUser, err
		}
	}
	return newUser, nil
}

//查

func IsUserExistByName(userName string) bool {

	var userExist = &model.User{}
	if err := dao.SqlSession.Model(&model.User{}).Where("name=?", userName).First(&userExist).Error; gorm.IsRecordNotFoundError(err) {
		//关注不存在
		return false
	}
	//关注存在
	return true
}

func IsUserExist(userName string, password string, login *model.User) error {
	if login == nil {
		return common.ErrorNullPointer
	}
	// 记录一下原始密码(用户登录的密码)
	originPassword := password
	// 新生成加密密码用于和查询到的密码比较
	newPassword := EncryptPassword([]byte(originPassword))
	dao.SqlSession.Where("name=? and password=?", userName, newPassword).First(login)
	if login.Model.ID == 0 {
		return common.ErrorFullPossibility
	}
	return nil
}

// GetUser 根据用户id获取用户信息
func GetUser(userId uint) (model.User, error) {
	//1.数据模型准备
	var user model.User
	//2.在users表中查对应user_id的user
	if err := dao.SqlSession.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetUserById 根据用户id获取用户信息，用于userInfo
func GetUserById(userId uint, user *model.User) error {
	if user == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("id=?", userId).First(user)
	return nil

}

// IsUserLegal 用户名和密码合法性检验
func IsUserLegal(userName string, passWord string) error {
	//1.用户名检验
	if userName == "" {
		return common.ErrorUserNameNull
	}
	if len(userName) > MaxUsernameLength {
		return common.ErrorUserNameExtend
	}
	//2.密码检验
	if passWord == "" {
		return common.ErrorPasswordNull
	}
	if len(passWord) > MaxPasswordLength || len(passWord) < MinPasswordLength {
		return common.ErrorPasswordLength
	}
	return nil
}

// EncryptPassword md5密码加密
func EncryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}
