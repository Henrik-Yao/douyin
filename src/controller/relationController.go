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

// RelationAction 关注/取消关注操作
func RelationAction(c *gin.Context) {
	//1.取数据
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	getToUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	getActionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	hostId := tokenStruct.UserId
	guestId := uint(getToUserId)
	actionType := uint(getActionType)

	//2.service层处理
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
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostId := tokenStruct.UserId

	//2.service层处理
	followingList, err := service.FollowingList(hostId)
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
			UserList: followingList,
		})
	}
}

// FollowerList 获取用户关注列表
func FollowerList(c *gin.Context) {

	//1.数据预处理
	getUserId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	hostId := uint(getUserId)

	//2.service层处理
	followingList, err := service.FollowerList(hostId)
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
			UserList: followingList,
		})
	}
}
