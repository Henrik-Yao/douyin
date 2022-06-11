package controller

import (
	"douyin/src/common"
	"douyin/src/middleware"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedVideo struct {
	Id            uint     `json:"id,omitempty"`
	Author        FeedUser `json:"author,omitempty"`
	PlayUrl       string   `json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount uint     `json:"favorite_count,omitempty"`
	CommentCount  uint     `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
	Title         string   `json:"title,omitempty"`
}
type FeedResponse struct {
	common.Response
	VideoList []FeedVideo `json:"video_list,omitempty"`
	NextTime  uint        `json:"next_time,omitempty"`
}
type FeedNoVideoResponse struct {
	common.Response
	NextTime uint `json:"next_time"`
}
type FeedUser struct {
	Id             uint   `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    uint   `json:"follow_count,omitempty"`
	FollowerCount  uint   `json:"follower_count,omitempty"`
	IsFollow       bool   `json:"is_follow,omitempty"`
	TotalFavorited uint   `json:"total_favorited"`
	FavoriteCount  uint   `json:"favorite_count"`
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
	lastTime, err := strconv.ParseInt(strLastTime, 10, 32)
	if err != nil {
		lastTime = 0
	}

	var feedVideoList []FeedVideo
	feedVideoList = make([]FeedVideo, 0)
	videoList, _ := service.FeedGet(lastTime)
	var newTime int64 = 0 //返回的视频的最久的一个的时间
	for _, x := range videoList {
		var tmp FeedVideo
		tmp.Id = x.ID
		tmp.PlayUrl = x.PlayUrl
		//tmp.Author = //依靠用户信息接口查询
		var user, err = service.GetUser(x.AuthorId)
		var feedUser FeedUser
		if err == nil { //用户存在
			feedUser.Id = user.ID
			feedUser.FollowerCount = user.FollowerCount
			feedUser.FollowCount = user.FollowCount
			feedUser.Name = user.Name
			//add
			feedUser.TotalFavorited = user.TotalFavorited
			feedUser.FavoriteCount = user.FavoriteCount
			feedUser.IsFollow = false
			if haveToken {
				// 查询是否关注
				tokenStruct, ok := middleware.CheckToken(strToken)
				if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
					var uid1 = tokenStruct.UserId //用户id
					var uid2 = x.AuthorId         //视频发布者id
					if service.IsFollowing(uid1, uid2) {
						feedUser.IsFollow = true
					}
				}
			}
		}
		tmp.Author = feedUser
		tmp.CommentCount = x.CommentCount
		tmp.CoverUrl = x.CoverUrl
		tmp.FavoriteCount = x.FavoriteCount
		tmp.IsFavorite = false
		if haveToken {
			//查询是否点赞过
			tokenStruct, ok := middleware.CheckToken(strToken)
			if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
				var uid = tokenStruct.UserId         //用户id
				var vid = x.ID                       // 视频id
				if service.CheckFavorite(uid, vid) { //有点赞记录
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
			NextTime:  uint(newTime),
		})
	} else {
		c.JSON(http.StatusOK, FeedNoVideoResponse{
			Response: common.Response{StatusCode: 0}, //成功
			NextTime: 0,                              //重新循环
		})
	}

}
