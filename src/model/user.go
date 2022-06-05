package model

import (
	"crypto/md5"
	"douyin/src/common"
	"douyin/src/dao"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

type User struct {
	gorm.Model
	Name          string `json:"name"`
	Password      string `json:"password"`
	FollowCount   int32  `json:"follow_count"`
	FollowerCount int32  `json:"follower_count"`
}

//userInfo 用户信息相关sql

type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once //单例
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

// QueryUserInfoById 根据用户ID查询用户信息
func (u *UserInfoDAO) QueryUserInfoById(userId int64, userinfo *User) error {
	if userinfo == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("id=?", userId).First(userinfo)
	//id为零值，说明sql执行失败
	if userinfo.Model.ID == 0 {
		return common.ErrorUserNotExit
	}
	return nil
}

// AddUserInfo 添加用户信息
func (u *UserInfoDAO) AddUserInfo(userinfo *User) error {
	if userinfo == nil {
		return common.ErrorNullPointer
	}
	return dao.SqlSession.Create(userinfo).Error
}

// IsUserExistById 根据用户ID判断用户是否存在
func (u *UserInfoDAO) IsUserExistById(id int64) bool {
	var userinfo User
	if err := dao.SqlSession.Where("id=?", id).Select("id").First(&userinfo).Error; err != nil {
		log.Println(err)
	}
	if userinfo.Model.ID == 0 {
		return false
	}
	return true
}

//userLogin

type UserLoginDAO struct {
}

var (
	userLoginDao  *UserLoginDAO
	userLoginOnce sync.Once
)

func NewUserLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = new(UserLoginDAO)
	})
	return userLoginDao
}

// QueryUserLogin 登录功能
func (u *UserLoginDAO) QueryUserLogin(username, password string, login *User) error {
	if login == nil {
		return common.ErrorNullPointer
	}
	originPassword := password // 记录一下原始密码(用户登录的密码)
	// 新生成加密密码用于和查询到的密码比较
	newPassword := EncryptPassword([]byte(originPassword))
	dao.SqlSession.Where("name=? and password=?", username, newPassword).First(login)
	if login.Model.ID == 0 {
		return common.ErrorFullPossibility
	}
	return nil
}

// IsUserExistByUsername 根据用户名检查用户是否存在
func (u *UserLoginDAO) IsUserExistByUsername(username string) bool { //是否有比较把返回值类型改为error
	var userLogin User
	dao.SqlSession.Where("username=?", username).First(&userLogin)
	if userLogin.Model.ID == 0 {
		return false
	}
	return true
}

const secret = "return11111"

// EncryptPassword md5密码加密
func EncryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}
