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
	r.GET("/getToken", func(c *gin.Context) {
		token, err := middleware.CreateToken(11111, "return")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	r.POST("/testToken", middleware.JwtMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "鉴权成功"})
	})

	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("/user")
		{
			userGroup.POST("/test", middleware.JwtMiddleware(), controller.CreateUser)
		}

		// publish路由组
		publishGroup := douyinGroup.Group("/publish")
		{
			publishGroup.POST("/action", controller.Publish)
			publishGroup.GET("/list", controller.PublishList)

		}
		// feed只有一层，不需要组了
		douyinGroup.GET("/feed/", controller.Feed)

		//// feed路由组
		//feedGroup := douyinGroup.Group("user")
		//{
		//	feedGroup.POST("/test", controller.CreateUser)
		//}
		//
		//// favorite路由组
		favoriteGroup := douyinGroup.Group("favorite")
		{
			favoriteGroup.POST("/action", controller.Favorite)
			favoriteGroup.GET("/list",controller.FavoriteList)
		}
		//
		//// comment路由组
		//commentGroup := douyinGroup.Group("user")
		//{
		//	commentGroup.POST("/test", controller.CreateUser)
		//}
		//
		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action", controller.RelationAction)
			relationGroup.GET("/follow/list", controller.FollowList)
			relationGroup.GET("/follower/list", controller.FollowerList)
		}
	}

	return r
}
