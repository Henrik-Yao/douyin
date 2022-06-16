package controller

import (
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	logging "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ReturnAuthor struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
type ReturnMyself struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
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
type ReturnVideo2 struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnMyself `json:"author"`
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
type VideoListResponse2 struct {
	common.Response
	VideoList []ReturnVideo2 `json:"video_list"`
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
	f, err := data.Open()
	if err != nil {
		err.Error()
	}
	playUrl, err := service.CosUpload(finalName, f)
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

	//直接传至云端，不用存储到本地
	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
	img := service.ExampleReadFrameAsJpeg(saveFile, 3) //获取第3帧封面
	//img, _ := jpeg.Decode(buf)//保存到本地时要用到
	//imgw, _ := os.Create(saveImage) //先创建，后写入
	//jpeg.Encode(imgw, img, &jpeg.Options{100})
	coverUrl, err := service.CosUpload(coverName, img)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//删除保存在本地中的视频
	err = os.Remove(saveFile) // ignore_security_alert
	if err != nil {
		logging.Info(err)
	}

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

	if GuestId == 0 || GuestId == HostId {
		//根据token-id查找用户
		getUser, err := service.GetUser(HostId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "Not find this person.",
			})
			c.Abort()
			return
		}

		returnMyself := ReturnMyself{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
		}
		//根据用户id查找 所有相关视频信息
		var videoList []model.Video
		videoList = service.GetVideoList(HostId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "null",
				},
				VideoList: nil,
			})
		} else { //需要展示的列表信息
			var returnVideoList2 []ReturnVideo2
			for i := 0; i < len(videoList); i++ {
				returnVideo2 := ReturnVideo2{
					VideoId:       videoList[i].ID,
					Author:        returnMyself,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList2 = append(returnVideoList2, returnVideo2)
			}
			c.JSON(http.StatusOK, VideoListResponse2{
				Response: common.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
				VideoList: returnVideoList2,
			})
		}
	} else {
		//根据传入id查找用户
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
			AuthorId:      getUser.ID,
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
}
