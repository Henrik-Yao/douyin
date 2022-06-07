package controller

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/model"
	"douyin/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CommentListResponse struct {
	common.Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment,omitempty"`
}

type UserResponse struct {
	ID            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint   `json:"follow_count,omitempty"`
	FollowerCount uint   `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type CommentResponse struct {
	ID         uint         `json:"id,omitempty"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
	User       UserResponse `json:"user,omitempty"`
}

func CommentAction(c *gin.Context) {

	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}

	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	// Unsupported type
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "Unsupported actionType",
		})
		c.Abort()
		return
	}

	if actionType == "1" { // post
		text := c.Query("comment_text")
		PostComment(c, userId, text, uint(videoId))
	} else if actionType == "2" { //delete
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 10)
		DeleteComment(c, uint(videoId), uint(commentId))
	}

}

func PostComment(c *gin.Context, userId uint, text string, videoId uint) {

	newComment := model.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: text,
	}

	// Post a comment and change the comment count (using database transaction)
	err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		if err := service.PostComment(newComment); err != nil {
			return err
		}
		if err := service.AddCommentCount(videoId); err != nil {
			return err
		}
		return nil
	})
	getUser, err2 := service.GetUser(userId)
	videoAuthor, err3 := service.GetVideoAuthor(videoId)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to post comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "post the comment successfully",
		},
		Comment: CommentResponse{
			ID:         newComment.ID,
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userId, videoAuthor),
			},
		},
	})
}

func DeleteComment(c *gin.Context, videoId uint, commentId uint) {

	// Remove a comment and reduce the comment count (using database transaction)
	err := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		if err := service.DeleteComment(commentId); err != nil {
			return err
		}
		if err := service.ReduceCommentCount(videoId); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to delete comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "delete the comment successfully",
	})
}

func CommentList(c *gin.Context) {

	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)
	commentList, err := service.GetCommentList(uint(videoId))

	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}

	var responseCommentList []CommentResponse
	for i := 0; i < len(commentList); i++ {
		getUser, err1 := service.GetUser(commentList[i].UserId)

		if err1 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseCommentList[i] = CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userId, commentList[i].ID),
			},
		}
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "Successfully obtained the comment list.",
		},
		CommentList: responseCommentList,
	})

}
