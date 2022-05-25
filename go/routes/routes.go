package routes

import (
	"douyin/go/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 主路由组
	douyinGroup := r.Group("douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("user")
		{
			userGroup.POST("/test", controller.CreateUser)
		}

		// publish路由组
		publishGroup := douyinGroup.Group("user")
		{
			publishGroup.POST("/test", controller.CreateUser)
		}

		// feed路由组
		feedGroup := douyinGroup.Group("user")
		{
			feedGroup.POST("/test", controller.CreateUser)
		}

		// favorite路由组
		favoriteGroup := douyinGroup.Group("user")
		{
			favoriteGroup.POST("/test", controller.CreateUser)
		}

		// comment路由组
		commentGroup := douyinGroup.Group("user")
		{
			commentGroup.POST("/test", controller.CreateUser)
		}

		// relation路由组
		relationGroup := douyinGroup.Group("user")
		{
			relationGroup.POST("/test", controller.CreateUser)
		}
	}

	return r
}
