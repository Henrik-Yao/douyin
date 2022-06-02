package controller

import (
	"douyin/go/dao"
	"douyin/go/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	model.Response
	UserList []model.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	//1.验证token
	//token := c.Query("token")

	//tokenStruck, ok := middleware.CheckToken(token)
	//if !ok {
	//	c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "token不正确"})
	//	c.Abort() //阻止执行
	//	return
	//}
	////token超时
	//if time.Now().Unix() > tokenStruck.ExpiresAt {
	//	c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "token过期"})
	//	c.Abort() //阻止执行
	//	return
	//}
	fmt.Println("token通过验证")

	//2.获取user_id, to_user_id 和action_type
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	to_user_id, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	action_type_, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	action_type := int32(action_type_)
	fmt.Println("user_id:", user_id)

	//err := service.PostRelationAction(user_id,to_user_id,action_type)
	//3.准备数据
	newrelation := model.Relation{UserId: user_id, ToUserId: to_user_id}

	//dao.SqlSession.AutoMigrate(&model.Relation{})         //模型关联到数据库表videos
	//4.action_type判断
	if action_type == 1 {
		//关注
		fmt.Println("关注")
		var getrelation = model.Relation{}
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", user_id, to_user_id).First(&getrelation).Error; gorm.IsRecordNotFoundError(err) {
			//找不到数据
			fmt.Println("关注不存在，新建关注")
			dao.SqlSession.Table("relations").Create(&newrelation)
		} else {
			//找到数据
			fmt.Println("关注已存在")
		}

	}
	if action_type == 2 {
		//取消关注
		fmt.Println("取消关注")
		var getrelation = model.Relation{}
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", user_id, to_user_id).First(&getrelation).Error; gorm.IsRecordNotFoundError(err) {
			//找不到数据
			fmt.Println("关注不存在，无需操作")
			//dao.SqlSession.Table("relations").Create(&newrelation)
		} else {
			//找到数据
			fmt.Println("关注已存在，取消关注")
			dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", user_id, to_user_id).Delete(&newrelation)
		}

	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	//新建follow表
	//新建fan表
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		//UserList: []model.User{DemoUserinfo},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		//UserList: []model.User{DemoUserinfo},
	})
}
