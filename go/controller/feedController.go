package controller

import (
	"douyin/go/common"
	"douyin/go/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	common.Response
	VideoList []model.FeedVideo `json:"video_list,omitempty"`
	NextTime  int64             `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  common.Response{StatusCode: 0}, //成功
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
