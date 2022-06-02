package controller

/*
说明：将token.go中的BindJSON()换成了ShouldBind().
*/

import (
	"douyin/go/dao"
	"douyin/go/middleware"
	"douyin/go/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

//上传视频方法
func Publish(c *gin.Context) {
	//1.验证token
	token := c.PostForm("token")
	tokenStruck, ok := middleware.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "token不正确"})
		c.Abort() //阻止执行
		return
	}
	//token超时
	if time.Now().Unix() > tokenStruck.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "token过期"})
		c.Abort() //阻止执行
		return
	}
	fmt.Println("token通过验证")

	//获取id 和 name
	user_id := tokenStruck.UserId
	//username := tokenStruck.UserName
	//fmt.Printf("%#v", user_id)
	//fmt.Printf("%#v", username)

	//2.接收传来的数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//2.视频存入本地，列表信息存入数据库
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user_id, filename)
	saveFile := filepath.Join("./videos/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + "--上传成功",
	})

	//3.建数据库表Videos
	dao.SqlSession.AutoMigrate(&model.UserLoginInfo{}) //如果已经被创建，则将其注释掉
	dao.SqlSession.AutoMigrate(&model.Video{})         //模型关联到数据库表videos

	title := c.PostForm("title") //获取传入的标题
	var userLogin model.UserLoginInfo
	dao.SqlSession.Table("user_login_infos").Where("user_id=?", user_id).First(&userLogin)
	//将对象Object序列化成json存储,试过了，后面获取video_list不可行
	//user, err := json.Marshal(userLogin)
	//if err != nil {
	//	fmt.Println("error:", err)
	//}

	video := model.VideoAll{
		UserId:        userLogin.UserId,
		Name:          userLogin.Name,
		FollowCount:   userLogin.FollowCount,
		FollowerCount: userLogin.FollowerCount,
		IsFollow:      userLogin.IsFollow,
		PlayUrl:       "../video/" + finalName,
		CoverUrl:      "......",
		FavoriteCount: 0, //需要从其他接口获取
		CommentCount:  0,
		IsFavorite:    false,
		Title:         title,
	}
	dao.SqlSession.Table("videos").Create(&video) //创建记录

}

//获取列表方法
func PublishList(c *gin.Context) {
	//1.鉴权token
	token := c.Query("token")
	tokenStruck, ok := middleware.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "token不正确"})
		c.Abort() //阻止执行
		return
	}
	//token超时
	if time.Now().Unix() > tokenStruck.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "token过期"})
		c.Abort() //阻止执行
		return
	}
	fmt.Println("token通过验证")

	//2.查询当前id用户的所有视频，返回页面
	user_id := c.Query("user_id")
	var video_list []model.Video
	dao.SqlSession.Table("videos").Where("user_id=?", user_id).Find(&video_list)
	if video_list == nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "No query found.",
			},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: video_list,
		})
	}
}
