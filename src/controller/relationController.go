package controller

import (
	"douyin/src/common"
	"douyin/src/middleware"
	"douyin/src/model"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReturnFollower struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type FollowingListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

type FollowerListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
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
	var err error
	var userList []model.User
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	hostId := tokenStruct.UserId
	guestId := uint(getGuestId)

	//2.判断查询类型
	if guestId == 0 {
		//查本人的关注表
		userList, err = service.FollowingList(hostId)
	} else {
		//查对方的关注表
		userList, err = service.FollowingList(guestId)
	}

	var ReturnFollowerList = make([]ReturnFollower, len(userList))
	for i, m := range userList {
		ReturnFollowerList[i].Id = m.ID
		ReturnFollowerList[i].Name = m.Name
		ReturnFollowerList[i].FollowCount = m.FollowCount
		ReturnFollowerList[i].FollowerCount = m.FollowerCount
		ReturnFollowerList[i].IsFollow = service.IsFollowing(hostId, m.ID)
	}

	//fmt.Println(ReturnFollowerList)
	//3.响应返回
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
			UserList: ReturnFollowerList,
		})
	}
}

// FollowerList 获取用户关注列表
func FollowerList(c *gin.Context) {

	//1.数据预处理
	var err error
	var userList []model.User
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	hostId := tokenStruct.UserId
	guestId := uint(getGuestId)

	//2.判断查询类型
	if guestId == 0 {
		//查本人的关注表
		userList, err = service.FollowerList(hostId)
	} else {
		//查对方的关注表
		userList, err = service.FollowerList(guestId)
	}

	var ReturnFollowerList = make([]ReturnFollower, len(userList))
	for i, m := range userList {
		ReturnFollowerList[i].Id = m.ID
		ReturnFollowerList[i].Name = m.Name
		ReturnFollowerList[i].FollowCount = m.FollowCount
		ReturnFollowerList[i].FollowerCount = m.FollowerCount
		ReturnFollowerList[i].IsFollow = service.IsFollowing(hostId, m.ID)
	}

	//3.处理
	if err != nil {
		c.JSON(http.StatusBadRequest, FollowerListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FollowerListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: ReturnFollowerList,
		})
	}
}
