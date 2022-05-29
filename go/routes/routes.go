package routes

import (
	"douyin/go/controller"
	"douyin/go/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 以下代码为鉴权中间件测试
	r.GET("getToken", func(c *gin.Context) {
		token, err := middleware.CreateToken(11111, "return")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	r.POST("testToken", middleware.JwtMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "鉴权成功"})
	})

	// 主路由组
	douyinGroup := r.Group("douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("user")
		{
			userGroup.POST("/test", controller.CreateUser)
		}

		//// publish路由组
		//publishGroup := douyinGroup.Group("user")
		//{
		//	publishGroup.POST("/test", controller.CreateUser)
		//}
		//
		//// feed路由组
		//feedGroup := douyinGroup.Group("user")
		//{
		//	feedGroup.POST("/test", controller.CreateUser)
		//}
		//
		//// favorite路由组
		//favoriteGroup := douyinGroup.Group("user")
		//{
		//	favoriteGroup.POST("/test", controller.CreateUser)
		//}
		//
		//// comment路由组
		//commentGroup := douyinGroup.Group("user")
		//{
		//	commentGroup.POST("/test", controller.CreateUser)
		//}
		//
		//// relation路由组
		//relationGroup := douyinGroup.Group("user")
		//{
		//	relationGroup.POST("/test", controller.CreateUser)
		//}
	}

	return r
}
