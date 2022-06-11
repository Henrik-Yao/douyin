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

// ReturnFollower 关注表与粉丝表共用的用户数据模型
type ReturnFollower struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// FollowingListResponse 关注表相应结构体
type FollowingListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

// FollowerListResponse 粉丝表相应结构体
type FollowerListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

// RelationAction 关注/取消关注操作
func RelationAction(c *gin.Context) {
	//1.取数据
	//1.1 从token中获取用户id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostId := tokenStruct.UserId
	//1.2 获取待关注的用户id
	getToUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	guestId := uint(getToUserId)
	//1.3 获取关注操作（关注1，取消关注2）
	getActionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	actionType := uint(getActionType)

	//2.自己关注/取消关注自己不合法
	if hostId == guestId {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "无法关注自己",
		})
		c.Abort()
		return
	}

	//3.service层进行关注/取消关注处理
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
	//1.1获取用户本人id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostId := tokenStruct.UserId
	//1.2获取其他用户id
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestId := uint(getGuestId)

	//2.判断查询类型，从数据库取用户列表
	var err error
	var userList []model.User
	if guestId == 0 {
		//若其他用户id为0，代表查本人的关注表
		userList, err = service.FollowingList(hostId)
	} else {
		//若其他用户id不为0，代表查对方的关注表
		userList, err = service.FollowingList(guestId)
	}

	//构造返回的数据
	var ReturnFollowerList = make([]ReturnFollower, len(userList))
	for i, m := range userList {
		ReturnFollowerList[i].Id = m.ID
		ReturnFollowerList[i].Name = m.Name
		ReturnFollowerList[i].FollowCount = m.FollowCount
		ReturnFollowerList[i].FollowerCount = m.FollowerCount
		ReturnFollowerList[i].IsFollow = service.IsFollowing(hostId, m.ID)
	}

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

// FollowerList 获取用户粉丝列表
func FollowerList(c *gin.Context) {

	//1.数据预处理
	//1.1获取用户本人id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostId := tokenStruct.UserId
	//1.2获取其他用户id
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestId := uint(getGuestId)

	//2.判断查询类型
	var err error
	var userList []model.User
	if guestId == 0 {
		//查本人的关注表
		userList, err = service.FollowerList(hostId)
	} else {
		//查对方的关注表
		userList, err = service.FollowerList(guestId)
	}

	//3.判断查询类型，从数据库取用户列表
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
