package routes

import (
	"douyin/src/controller"
	"douyin/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("/user")
		{
			userGroup.POST("/test", middleware.JwtMiddleware(), controller.CreateUser)
			userGroup.GET("/user/", middleware.JwtMiddleware(), controller.UserInfo)
			userGroup.POST("/user/login/", middleware.JwtMiddleware(), controller.UserLogin)
			userGroup.POST("/user/register/", middleware.JwtMiddleware(), controller.UserRegister)
		}

		// publish路由组
		publishGroup := douyinGroup.Group("/publish")
		{
			publishGroup.POST("/action", controller.Publish) //提交文件，不用中间件鉴权
			publishGroup.GET("/list", middleware.JwtMiddleware(), controller.PublishList)

		}
		// feed只有一层，不需要组了
		douyinGroup.GET("/feed/", controller.Feed)

		favoriteGroup := douyinGroup.Group("favorite")
		{
			favoriteGroup.POST("/action", middleware.JwtMiddleware(), controller.Favorite)
			favoriteGroup.GET("/list", middleware.JwtMiddleware(), controller.FavoriteList)
		}
		//
		//// comment路由组
		commentGroup := douyinGroup.Group("/comment")
		{
			commentGroup.POST("/action", middleware.JwtMiddleware(), controller.CommentAction)
			commentGroup.GET("/list", middleware.JwtMiddleware(), controller.CommentList)
		}
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
