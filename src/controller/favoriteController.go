package controller

import (
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteAuthor struct { //从user中获取,getUser函数
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"` //从following或follower中获取
}

type FavoriteVideo struct { //从video中获取
	Id            uint           `json:"id,omitempty"`
	Author        FavoriteAuthor `json:"author,omitempty"`
	PlayUrl       string         `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string         `json:"cover_url,omitempty"`
	FavoriteCount uint           `json:"favorite_count,omitempty"`
	CommentCount  uint           `json:"comment_count,omitempty"`
	IsFavorite    bool           `json:"is_favorite,omitempty"` //true
	Title         string         `json:"title,omitempty"`
}

type FavoriteListResponse struct {
	common.Response
	VideoList []FavoriteVideo `json:"video_list,omitempty"`
}

// Favorite 点赞视频方法
func Favorite(c *gin.Context) {
	//参数绑定
	//user_id获取
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	//参数获取
	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseUint(actionTypeStr, 10, 10)
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	//函数调用及响应
	err := service.FavoriteAction(userId, uint(videoId), uint(actionType))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "操作成功！",
		})
	}
}

// FavoriteList 获取列表方法
func FavoriteList(c *gin.Context) {
	//user_id获取
	getUserId, _ := c.Get("user_id")
	var userIdHost uint
	if v, ok := getUserId.(uint); ok {
		userIdHost = v
	}
	userIdStr := c.Query("user_id") //自己id或别人id
	userId, _ := strconv.ParseUint(userIdStr, 10, 10)
	userIdNew := uint(userId)
	if userIdNew == 0 {
		userIdNew = userIdHost
	}

	//函数调用及响应
	videoList, err := service.FavoriteList(userIdNew)
	videoListNew := make([]FavoriteVideo, 0)
	for _, m := range videoList {
		var author = FavoriteAuthor{}
		var getAuthor = model.User{}
		getAuthor, err := service.GetUser(m.AuthorId) //参数类型、错误处理
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "找不到作者！",
			})
			c.Abort()
			return
		}
		//isfollowing
		isfollowing := service.IsFollowing(userIdHost, m.AuthorId) //参数类型、错误处理
		//isfavorite
		isfavorite := service.CheckFavorite(userIdHost, m.ID)
		//作者信息
		author.Id = getAuthor.ID
		author.Name = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = isfollowing
		//组装
		var video = FavoriteVideo{}
		video.Id = m.ID //类型转换
		video.Author = author
		video.PlayUrl = m.PlayUrl
		video.CoverUrl = m.CoverUrl
		video.FavoriteCount = m.FavoriteCount
		video.CommentCount = m.CommentCount
		video.IsFavorite = isfavorite
		video.Title = m.Title

		videoListNew = append(videoListNew, video)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, FavoriteListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			VideoList: videoListNew,
		})
	}
}
