package controller

import (
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FavoriteListResponse struct {
	common.Response
	VideoList []model.Video `json:"video_list"`
}

//点赞视频方法
func Favorite(c *gin.Context) {
	//验证token
	//参数绑定
	var favoritereq model.FavoriteRequest
	c.BindJSON(&favoritereq)

	fmt.Println("token通过验证")

	//函数调用及响应
	err := service.FavoriteAction(&favoritereq)
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

//获取列表方法
func FavoriteList(c *gin.Context) {
	//鉴权token
	userId := c.Query("user_id")

	//函数调用及响应
	videoList, err := service.FavoriteList(userId)
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
			VideoList: videoList,
		})
	}
}
