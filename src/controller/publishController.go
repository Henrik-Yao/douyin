package controller

import (
	"context"
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	logging "github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type ReturnAuthor struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type ReturnVideo struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnAuthor `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

type VideoListResponse struct {
	common.Response
	VideoList []ReturnVideo `json:"video_list"`
}

const (
	COS_BUCKET_NAME = "dong"
	COS_APP_ID      = "1305843950"
	COS_REGION      = "ap-nanjing"
	COS_SECRET_ID   = "AKIDa0B5j6C1ZMvDG4brqZ1B5i5BzXprc6KH"
	COS_SECRET_KEY  = "h7G9l6AxTigNozuYuzoMfjX4NREl1KNA"
	COS_URL_FORMAT  = "http://%s-%s.cos.%s.myqcloud.com" // 此项固定
)

func CosUpload(fileName string, path string) (string, error) {
	u, _ := url.Parse(fmt.Sprintf(COS_URL_FORMAT, COS_BUCKET_NAME, COS_APP_ID, COS_REGION))
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  COS_SECRET_ID,
			SecretKey: COS_SECRET_KEY,
		},
	})

	//path为本地的保存路径
	_, err := client.Object.PutFromFile(context.Background(), fileName, path, nil)
	if err != nil {
		panic(err)
	}
	return "https://dong-1305843950.cos.ap-nanjing.myqcloud.com/" + fileName, nil
}

func Publish(c *gin.Context) { //上传视频方法
	//1.中间件验证token后，获取userId
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	//2.接收请求参数信息
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//3.返回至前端页面的展示信息
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, fileName)
	//先存储到本地文件夹，再保存到云端，获取封面后最后删除
	saveFile := filepath.Join("../videos/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//从本地上传到云端，并获取云端地址
	playUrl, err := CosUpload(finalName, saveFile)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + "--uploaded successfully",
	})

	var coverUrl string
	coverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"

	////将封面上传到云端并返回url
	//coverUrl, err := CosUpload(coverName, coverName)
	//if err != nil {
	//	c.JSON(http.StatusOK, common.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}

	//删除本地public中的视频
	err = os.Remove(saveFile)
	if err != nil {
		logging.Info(err)
	}
	//删除本地封面
	//err = os.Remove(coverPath)
	//if err != nil {
	//	logging.Info(err)
	//}

	//4.保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	video := model.Video{
		Model:         gorm.Model{},
		AuthorId:      userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	service.CreateVideo(&video)
}

func PublishList(c *gin.Context) { //获取列表的方法
	//1.中间件鉴权token
	getHostId, _ := c.Get("user_id")
	var HostId uint
	if v, ok := getHostId.(uint); ok {
		HostId = v
	}
	//2.查询要查看用户的id的所有视频，返回页面
	getGuestId := c.Query("user_id")
	id, _ := strconv.Atoi(getGuestId)
	GuestId := uint(id)

	//根据用户id查找用户
	getUser, err := service.GetUser(GuestId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  "Not find this person.",
		})
		c.Abort()
		return
	}

	returnAuthor := ReturnAuthor{
		AuthorId:      GuestId,
		Name:          getUser.Name,
		FollowCount:   getUser.FollowCount,
		FollowerCount: getUser.FollowerCount,
		IsFollow:      service.IsFollowing(HostId, GuestId),
	}
	//根据用户id查找 所有相关视频信息
	var videoList []model.Video
	videoList = service.GetVideoList(GuestId)
	if len(videoList) == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "null",
			},
			VideoList: nil,
		})
	} else { //需要展示的列表信息
		var returnVideoList []ReturnVideo
		for i := 0; i < len(videoList); i++ {
			returnVideo := ReturnVideo{
				VideoId:       videoList[i].ID,
				Author:        returnAuthor,
				PlayUrl:       videoList[i].PlayUrl,
				CoverUrl:      videoList[i].CoverUrl,
				FavoriteCount: videoList[i].FavoriteCount,
				CommentCount:  videoList[i].CommentCount,
				IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
				Title:         videoList[i].Title,
			}
			returnVideoList = append(returnVideoList, returnVideo)
		}
		c.JSON(http.StatusOK, VideoListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: returnVideoList,
		})
	}
}
