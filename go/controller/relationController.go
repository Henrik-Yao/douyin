package controller

import (
	"douyin/go/dao"
	"douyin/go/model"
	"douyin/go/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	model.Response
	UserList []follower `json:"user_list"`
}

type follower struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	//1.预处理
	token := c.Query("token")
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	relationreq := model.RelationRequest{
		UserId:     user_id,
		Token:      token,
		ToUserId:   to_user_id,
		ActionType: int32(action_type),
	}
	//2.token鉴权
	//token := realtionreq.Token

	//3.service层处理
	err := service.RelationAction(&relationreq)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "操作成功！",
		})
	}

}

// FollowList
func FollowList(c *gin.Context) {

	//1.验证token

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//2.先查relation表
	var getrelation []follower
	dao.SqlSession.Table("user_login_infos").Joins("left join relations on user_login_infos.user_id = relations.to_user_id").
		Where("relations.user_id=?", user_id).Scan(&getrelation)
	//is_fllow逻辑

	fmt.Println(getrelation)
	if getrelation == nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "No query found.",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			UserList: getrelation,
		})
	}
}

// FollowerList
func FollowerList(c *gin.Context) {
	//1.验证token

	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//2.先查relation表
	var getrelation []follower
	dao.SqlSession.Table("user_login_infos").Joins("left join relations on user_login_infos.user_id = relations.user_id").
		Where("relations.to_user_id=?", user_id).Scan(&getrelation)
	//is_fllow逻辑需要增加

	fmt.Println(getrelation)
	if getrelation == nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "No query found.",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			UserList: getrelation,
		})
	}
}
