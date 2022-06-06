package controller

import (
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)



type UserRegisterResponse struct {

	common.Response
	service.UserResponse
}

func UserRegister(c *gin.Context) {
	//1.参数提取
	username := c.Query("username")
	password := c.Query("password")

	//2.service层处理
	registerResponse, err := service.UserRegister(username, password)

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
		Response:     common.Response{StatusCode: 0},
		UserResponse: registerResponse,
	})
	return
}

type UserLoginResponse struct {
	common.Response
	service.UserResponse
}

func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := service.UserLogin(username, password)

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
		Response:     common.Response{StatusCode: 0},
		UserResponse: userLoginResponse,
	})
}

type UserInfoResponse struct {
	common.Response
	service.UserInfoQueryResponse
}

func UserInfo(c *gin.Context) {

	//根据user_id查询
	rawId := c.Query("user_id")
	userInfoResponse, err := service.UserInfo(rawId)

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
		Response:              common.Response{StatusCode: 0},
		UserInfoQueryResponse: userInfoResponse,
	})

}
