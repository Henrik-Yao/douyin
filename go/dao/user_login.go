package dao

import (
	"crypto/md5"
	"douyin/go/model"
	"encoding/hex"
	"sync"
)

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
func (u *UserLoginDAO) QueryUserLogin(username, password string, login *model.UserLogin) error {
	if login == nil {
		return ErrorNullPointer
	}
	originPassword := password // 记录一下原始密码(用户登录的密码)
	// 新生成加密密码用于和查询到的密码比较
	newPassword := EncryptPassword([]byte(originPassword))
	SqlSession.Where("username=? and password=?", username, newPassword).First(login)
	if login.Id == 0 {
		return ErrorFullPossibility
	}
	return nil
}

// IsUserExistByUsername 根据用户名检查用户是否存在
func (u *UserLoginDAO) IsUserExistByUsername(username string) bool { //是否有比较把返回值类型改为error
	var userLogin model.UserLogin
	SqlSession.Where("username=?", username).First(&userLogin)
	if userLogin.Id == 0 {
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
