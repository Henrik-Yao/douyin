package controller

import (
	"douyin/src/common"
	"douyin/src/middleware"
	"douyin/src/model"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
type FeedNoVideoResponse struct {
	common.Response
	NextTime int64 `json:"next_time"`
}

func Feed(c *gin.Context) {

	strToken := c.Query("token")
	var haveToken bool
	if strToken == "" {
		haveToken = false
	} else {
		haveToken = true
	}
	var strLastTime = c.Query("latest_time")
	lastTime, err := strconv.ParseInt(strLastTime, 10, 64)
	if err != nil {
		lastTime = 0
	}

	fmt.Println(lastTime)
	var feedVideoList []FeedVideo
	feedVideoList = make([]FeedVideo, 0)
	videoList, _ := service.FeedGet(lastTime)
	var newTime int64 = 0 //返回的视频的最久的一个的时间
	for _, x := range videoList {
		var tmp FeedVideo
		tmp.Id = int64(x.ID)
		tmp.PlayUrl = x.PlayUrl
		//tmp.Author = //依靠用户信息接口查询
		tmp.CommentCount = x.CommentCount
		tmp.CoverUrl = x.CoverUrl
		tmp.FavoriteCount = x.FavoriteCount
		tmp.IsFavorite = false
		if haveToken {
			//查询是否点赞过
			tokenStruct, ok := middleware.CheckToken(strToken)
			if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
				var uid = tokenStruct.UserId              //用户id
				var vid = x.ID                            // 视频id
				if service.CheckFavorite(uid, int(vid)) { //有点赞记录
					tmp.IsFavorite = true
				}
			}
		}
		tmp.Title = x.Title
		feedVideoList = append(feedVideoList, tmp)
		newTime = x.CreatedAt.Unix()
	}
	if len(feedVideoList) > 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  common.Response{StatusCode: 0}, //成功
			VideoList: feedVideoList,
			NextTime:  newTime,
		})
	} else {
		c.JSON(http.StatusOK, FeedNoVideoResponse{
			Response: common.Response{StatusCode: 0}, //成功
			NextTime: 0,                              //重新循环
		})
	}

}
