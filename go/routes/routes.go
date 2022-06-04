package routes

import (
	"douyin/go/controller"
	"douyin/go/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

		fmt.Println("-----", c.PostForm(string("id")))
		c.JSON(http.StatusOK, gin.H{"msg": "鉴权成功"})
	})
	r.POST("/testToken2", func(c *gin.Context) {
		s := c.Query("token")
		fmt.Println(s)
		fmt.Println("-----", c.PostForm("token"))
		c.JSON(http.StatusOK, gin.H{"msg": "鉴权成功"})
	})

	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// user路由组
		userGroup := douyinGroup.Group("/user")
		{
			userGroup.POST("/test", middleware.JwtMiddleware(), controller.CreateUser)
			userGroup.GET("/user/", middleware.JwtMiddleware(), controller.UserInfoHandler)
			userGroup.POST("/user/login/", middleware.JwtMiddleware(), controller.UserLoginHandler)
			userGroup.POST("/user/register/", middleware.JwtMiddleware(), controller.UserRegisterHandler)
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
			favoriteGroup.GET("/list", controller.FavoriteList)
		}
		//
		//// comment路由组
		commentGroup := douyinGroup.Group("/comment")
		{
			commentGroup.POST("/action", controller.CommentAction)
			commentGroup.GET("/list", controller.CommentList)
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
