package controller

import (
	"douyin/src/common"
	"douyin/src/middleware"
	"douyin/src/model"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserIdTokenResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserRegisterResponse struct {
	common.Response
	UserIdTokenResponse
}

// UserRegister 用户注册主函数
func UserRegister(c *gin.Context) {
	//1.参数提取
	username := c.Query("username")
	password := c.Query("password")

	//2.service层处理
	registerResponse, err := UserRegisterService(username, password)
	// UserRegister 用户注册

	//3.返回响应
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response:            common.Response{StatusCode: 0},
		UserIdTokenResponse: registerResponse,
	})
	return
}

// UserRegisterService 用户注册用户登录处理函数
func UserRegisterService(userName string, passWord string) (UserIdTokenResponse, error) {

	//0.数据准备
	var userResponse = UserIdTokenResponse{}

	//1.合法性检验
	err := service.IsUserLegal(userName, passWord)
	if err != nil {
		return userResponse, err
	}

	//2.新建用户
	newUser, err := service.CreateRegisterUser(userName, passWord)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserIdTokenResponse{
		UserId: newUser.ID,
		Token:  token,
	}
	return userResponse, nil
}

type UserLoginResponse struct {
	common.Response
	UserIdTokenResponse
}

// UserLogin 用户登录主函数
func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := UserLoginService(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:            common.Response{StatusCode: 0},
		UserIdTokenResponse: userLoginResponse,
	})
}

// UserLoginService 用户登录处理函数
func UserLoginService(userName string, passWord string) (UserIdTokenResponse, error) {

	//0.数据准备
	var userResponse = UserIdTokenResponse{}

	//1.合法性检验
	err := service.IsUserLegal(userName, passWord)
	if err != nil {
		return userResponse, err
	}

	//2.检查用户是否存在
	var login model.User
	err = service.IsUserExist(userName, passWord, &login)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(login.Model.ID, login.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserIdTokenResponse{
		UserId: login.Model.ID,
		Token:  token,
	}
	return userResponse, nil
}

// UserInfoQueryResponse 用户信息返回的结构体
type UserInfoQueryResponse struct {
	UserId        uint   `json:"user_id"`
	UserName      string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserInfoResponse struct {
	common.Response
	UserList UserInfoQueryResponse `json:"user"`
}

// UserInfo 用户信息主函数
func UserInfo(c *gin.Context) {
	//根据user_id查询
	rawId := c.Query("user_id")
	userInfoResponse, err := UserInfoService(rawId)

	//根据token获得当前用户的userid
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)
	hostId := tokenStruct.UserId
	userInfoResponse.IsFollow = service.CheckIsFollow(rawId, hostId)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserList: userInfoResponse,
	})

}

// UserInfoService 用户信息处理函数
func UserInfoService(rawId string) (UserInfoQueryResponse, error) {
	//0.数据准备
	var userInfoQueryResponse = UserInfoQueryResponse{}
	userId, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		return userInfoQueryResponse, err
	}

	//1.获取用户信息
	var user model.User
	err = service.GetUserById(uint(userId), &user)
	if err != nil {
		return userInfoQueryResponse, err
	}

	userInfoQueryResponse = UserInfoQueryResponse{
		UserId:        user.Model.ID,
		UserName:      user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}
	return userInfoQueryResponse, nil
}
