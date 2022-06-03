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
	"github.com/jinzhu/gorm"
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
	//fmt.Println("token通过验证")

	//获取id
	user_id := tokenStruck.UserId

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

	//3.建数据库表
	//dao.SqlSession.AutoMigrate(&model.UserInfo{})
	dao.SqlSession.AutoMigrate(&model.VideoNoAuthor{}) //模型关联到数据库表

	////保存视频-用户关系信息(可以不用，多此一举)
	//video_id = video_id + 1
	//videoUser := model.VideoUser{
	//	Model:   gorm.Model{},
	//	UserId:  int64(user_id),
	//	VideoId: video_id,
	//}
	//dao.SqlSession.Table("video_users").Create(&videoUser)

	title := c.PostForm("title") //获取传入的标题
	//保存视频相关信息
	videoNoAuthor := model.VideoNoAuthor{
		Model:         gorm.Model{},
		UserId:        int64(user_id),
		PlayUrl:       "../videos/" + finalName,
		CoverUrl:      "......",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         title,
	}
	dao.SqlSession.Table("video_no_authors").Create(&videoNoAuthor)

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

	//2.查询当前id用户的所有视频，返回页面
	user_id := c.Query("user_id")
	//根据id查找用户
	var userInfo model.UserInfo
	dao.SqlSession.Table("user_infos").Where("user_id=?", user_id).First(&userInfo)
	loginUser := model.LoginUser{
		UserId:        userInfo.UserId,
		Name:          userInfo.Name,
		FollowCount:   userInfo.FollowCount,
		FollowerCount: userInfo.FollowerCount,
		IsFollow:      userInfo.IsFollow,
	}
	//根据用户id查找 所有相关视频信息
	var videoInfoList []model.VideoNoAuthor
	dao.SqlSession.Table("video_no_authors").Where("user_id=?", user_id).Find(&videoInfoList)

	if len(videoInfoList) == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "No query found.",
			},
			VideoList: nil,
		})
	} else {
		//需要展示的列表信息
		var videoList []model.Video
		for i := 0; i < len(videoInfoList); i++ {
			video := model.Video{
				VideoId:       int64(videoInfoList[i].ID),
				Author:        loginUser,
				PlayUrl:       videoInfoList[i].PlayUrl,
				CoverUrl:      videoInfoList[i].CoverUrl,
				FavoriteCount: videoInfoList[i].FavoriteCount,
				CommentCount:  videoInfoList[i].CommentCount,
				IsFavorite:    videoInfoList[i].IsFavorite,
				Title:         videoInfoList[i].Title,
			}
			videoList = append(videoList, video)
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: videoList,
		})
	}
}
