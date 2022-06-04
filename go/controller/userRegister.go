package controller

import (
	"douyin/go/common"
	"douyin/go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterResponse struct {
	common.Response
	*service.UserLoginResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	registerResponse, err := service.PostUserLogin(username, password)

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
		Response:          common.Response{StatusCode: 0},
		UserLoginResponse: registerResponse,
	})
}
