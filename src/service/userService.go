package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/middleware"
	"douyin/src/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	MinPasswordLength = 6  //密码最小长度
)

// UserResponse 用户注册登录返回的结构体
type UserResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

// UserInfoQueryResponse 用户信息返回的结构体
type UserInfoQueryResponse struct {
	UserId        uint   `json:"user_id"`
	UserName      string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
}

//增

// CreateUser 新建用户
func CreateUser(user *model.User) (err error) {
	if err = dao.SqlSession.Create(user).Error; err != nil {
		return err
	}
	return
}

func CreateRegisterUser(userName string, passWord string) (model.User, error) {
	//1.Following数据模型准备
	newUser := model.User{
		Name:     userName,
		Password: passWord,
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
	dao.SqlSession.Where("name=? and password=?", userName, password).First(login)
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

// UserRegister 用户注册
func UserRegister(userName string, passWord string) (UserResponse, error) {

	//0.数据准备
	var userResponse = UserResponse{}

	//1.合法性检验
	err := IsUserLegal(userName, passWord)
	if err != nil {
		return userResponse, err
	}

	//2.新建用户
	newUser, err := CreateRegisterUser(userName, passWord)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserResponse{
		UserId: newUser.ID,
		Token:  token,
	}
	return userResponse, nil
}

// UserLogin 用户登录
func UserLogin(userName string, passWord string) (UserResponse, error) {

	//0.数据准备
	var userResponse = UserResponse{}

	//1.合法性检验
	err := IsUserLegal(userName, passWord)
	if err != nil {
		return userResponse, err
	}

	//2.检查用户是否存在
	var login model.User
	err = IsUserExist(userName, passWord, &login)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(login.Model.ID, login.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserResponse{
		UserId: login.Model.ID,
		Token:  token,
	}
	return userResponse, nil
}

// UserInfo 用户信息
func UserInfo(rawId string) (UserInfoQueryResponse, error) {
	//0.数据准备
	var userInfoQueryResponse = UserInfoQueryResponse{}
	userId, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		return userInfoQueryResponse, err
	}

	//1.获取用户信息
	var user model.User
	err = GetUserById(uint(userId), &user)
	if err != nil {
		return userInfoQueryResponse, err
	}

	userInfoQueryResponse = UserInfoQueryResponse{
		UserId:        user.Model.ID,
		UserName:      user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
	}
	return userInfoQueryResponse, nil
}
