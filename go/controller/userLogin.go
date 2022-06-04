package controller

import (
	"douyin/go/model"
	"douyin/go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	model.Response
	*service.UserLoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := service.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:          model.Response{StatusCode: 0},
		UserLoginResponse: userLoginResponse,
	})
}
