package controller

import (
	"douyin/src/common"
	"douyin/src/middleware"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FollowingListResponse struct {
	common.Response
	UserList []service.Follower `json:"user_list"`
}

type FollowerListResponse struct {
	common.Response
	UserList []service.Follower `json:"user_list"`
}

// 关注/取消关注操作
func FollowAction(c *gin.Context) {
	//1.取数据
	user_id_, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	to_user_id_, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type_, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	hostId := int32(user_id_)
	guestId := int32(to_user_id_)
	actionType := int32(action_type_)

	//2.token鉴权
	middleware.JwtMiddleware()

	//3.actionType判断
	err := service.FollowAction(hostId, guestId, actionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "关注/取消关注成功！",
		})
	}
}

// FollowList 获取用户关注列表
func FollowList(c *gin.Context) {

	//1.数据预处理
	user_id_, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	hostId := int32(user_id_)

	//2.token鉴权
	middleware.JwtMiddleware()

	//3.service层处理
	followinglist, err := service.FollowingList(hostId)
	if err != nil {
		c.JSON(http.StatusBadRequest, FollowingListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FollowingListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: followinglist,
		})
	}
}

// FollowerList 获取用户关注列表
func FollowerList(c *gin.Context) {

	//1.数据预处理
	user_id_, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	hostId := int32(user_id_)

	//2.token鉴权
	middleware.JwtMiddleware()

	//3.service层处理
	followinglist, err := service.FollowerList(hostId)
	if err != nil {
		c.JSON(http.StatusBadRequest, FollowingListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FollowingListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: followinglist,
		})
	}
}
