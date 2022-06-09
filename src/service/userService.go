package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	MinPasswordLength = 6  //密码最小长度
)

//增

func CreateRegisterUser(userName string, passWord string) (model.User, error) {
	//1.Following数据模型准备
	newPassword, _ := HashAndSalt(passWord)
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
			//fmt.Println(err)
			panic(err)
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
	dao.SqlSession.Where("name=?", userName).First(login)
	if !ComparePasswords(login.Password, password) {
		return common.ErrorPasswordFalse
	}
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

// HashAndSalt 加密密码
func HashAndSalt(pwdStr string) (pwdHash string, err error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	pwdHash = string(hash)
	return
}

// ComparePasswords 验证密码
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
