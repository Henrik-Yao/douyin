package controller

import (
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedVideo struct {
	Id            int64      `json:"id,omitempty"`
	Author        model.User `json:"author,omitempty"`
	PlayUrl       string     `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string     `json:"cover_url,omitempty"`
	FavoriteCount int64      `json:"favorite_count,omitempty"`
	CommentCount  int64      `json:"comment_count,omitempty"`
	IsFavorite    bool       `json:"is_favorite,omitempty"`
	Title         string     `json:"title,omitempty"`
}
type FeedResponse struct {
	common.Response
	VideoList []FeedVideo `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {

	str := c.Query("token")
	if str == "" {
		fmt.Println("no token")

	} else {
		fmt.Println("token=", str)
	}
	var feedVideoList []FeedVideo
	feedVideoList = make([]FeedVideo, 0)
	videoList, _ := service.FeedGet()
	for _, x := range videoList {
		var tmp FeedVideo
		tmp.Id = int64(x.ID)
		tmp.PlayUrl = x.PlayUrl
		//tmp.Author = //依靠信息查询用户信息
		tmp.CommentCount = x.CommentCount
		tmp.CoverUrl = x.CoverUrl
		tmp.FavoriteCount = x.FavoriteCount
		tmp.IsFavorite = false
		tmp.Title = x.Title
		feedVideoList = append(feedVideoList, tmp)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  common.Response{StatusCode: 0}, //成功
		VideoList: feedVideoList,
		NextTime:  time.Now().Unix(),
	})
}
