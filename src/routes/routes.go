package routes

import (
	"douyin/src/controller"
	"douyin/src/middleware"
	"fmt"
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

		fmt.Println("-----", c.PostForm(string("id")))
		c.JSON(http.StatusOK, gin.H{"msg": "鉴权成功"})
	})
	r.POST("/testToken2", middleware.JwtMiddleware(), func(c *gin.Context) {
		s, _ := c.Get("user_id")
		var s2 int
		if v, ok := s.(int); ok {
			s2 = v
		}
		fmt.Println(s2)
		//fmt.Println("-----", s)
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
			commentGroup.POST("/action", controller.CommentAction)
			commentGroup.GET("/list", controller.CommentList)
		}
		//
		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action", controller.FollowAction)
			relationGroup.GET("/follow/list", controller.FollowList)
			relationGroup.GET("/follower/list", controller.FollowerList)
		}
	}

	return r
}
