package controller

import (
	"douyin/go/dao"
	"douyin/go/middleware"
	"douyin/go/model"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment CommentResponse `json:"comment,omitempty"`
}

type CommentResponse struct {
	model.Comment
	User model.UserLoginInfo `json:"user,omitempty"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	tokenStruct, ok := middleware.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "User doesn't exist"})
		c.Abort()
		return
	}
	if time.Now().Unix() > tokenStruct.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "Token has expired"})
		c.Abort()
		return
	}

	// Unsupported type
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "Unsupported actionType"})
		c.Abort()
		return
	}

	if actionType == "1" { // post
		text := c.Query("comment_text")
		PostComment(c, tokenStruct.UserId, text, videoId)
	} else if actionType == "2" { //delete
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
		DeleteComment(c, videoId, commentId)
	}

}

func PostComment(c *gin.Context, userId int, text string, videoId int64) {

	// Generate random numbers as id (self-increment id is difficult to obtain)
	var randomId int64
	for {
		rand.Seed(time.Now().UnixNano())
		partRand := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000)
		randomIdStr := strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(int64(partRand), 10)
		randomId, _ = strconv.ParseInt(randomIdStr, 10, 64)
		var userExist model.UserLoginInfo
		dao.SqlSession.Table("user_login_infos").Where("user_id=?", randomId).Find(&userExist)
		if userExist == (model.UserLoginInfo{}) {
			rand.Seed(time.Now().UnixNano())
			partRand := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000)
			randomIdStr := strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(int64(partRand), 10)
			randomId, _ = strconv.ParseInt(randomIdStr, 10, 64)
		} else {
			break
		}
	}

	newComment := model.Comment{
		Id:         randomId,
		UserId:     int64(userId),
		Content:    text,
		CreateDate: time.Now().String(),
		VideoId:    videoId,
		IsDeleted:  false,
	}

	dao.SqlSession.AutoMigrate(&model.Comment{})

	err := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		// Add a comment record
		if err := dao.SqlSession.Table("comment").Create(&newComment).Error; err != nil {
			return err
		}
		// Change the number of video comments
		dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1"))
		return nil
	})
	if err != nil {
		return
	}

	var getUser model.UserLoginInfo
	dao.SqlSession.Table("user_login_infos").Where("user_id=?", userId).Find(&getUser)
	currUser := model.UserLoginInfo{
		UserId:        getUser.UserId,
		Name:          getUser.Name,
		FollowCount:   getUser.FollowCount,
		FollowerCount: getUser.FollowerCount,
		IsFollow:      getUser.IsFollow,
	}

	c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
		Comment: CommentResponse{
			newComment,
			currUser,
		}})
}

func DeleteComment(c *gin.Context, videoId int64, commentId int64) {
	err := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		// Modify a field that indicates whether it has been deleted
		if err := dao.SqlSession.Table("comment").Where("id = ?", commentId).Update("is_deleted", true).Error; err != nil {
			return err
		}
		// Change the number of video comments
		dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1"))
		return nil
	})
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "Comments have been deleted successfully",
	})
}

func CommentList(c *gin.Context) {
	token := c.Query("token")

	tokenStruct, ok := middleware.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "User doesn't exist"})
		c.Abort()
		return
	}

	if time.Now().Unix() > tokenStruct.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "Token has expired"})
		c.Abort()
		return
	}

	videoId := c.Query("video_id")
	var commentList []model.Comment
	dao.SqlSession.Table("comment").Where("video_id=? and is_deleted = false", videoId).Find(&commentList)

	var responseCommentList []CommentResponse
	for i := 0; i < len(commentList); i++ {
		var getUser model.UserLoginInfo
		dao.SqlSession.Table("user_login_infos").Where("user_id=?", commentList[i].UserId).Find(&getUser)
		responseCommentList[i] = CommentResponse{
			commentList[i],
			getUser,
		}
	}

	if commentList == nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "No query found.",
			},
			CommentList: nil,
		})
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			CommentList: responseCommentList,
		})
	}
}