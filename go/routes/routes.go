package routes

import (
	"douyin/go/controller"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	// 主路由组
	douyinGroup := r.Group("douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("user")
		{

			userGroup.POST("/test", controller.CreateUser)
		}

	}

	return r
}
