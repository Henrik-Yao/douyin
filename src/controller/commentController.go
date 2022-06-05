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

type UserRespoonse struct {
	ID            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int32  `json:"follow_count,omitempty"`
	FollowerCount int32  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
type CommentResponse struct {
	ID         int64         `json:"id,omitempty"`
	Content    string        `json:"content,omitempty"`
	CreateDate string        `json:"create_date,omitempty"`
	User       UserRespoonse `json:"user,omitempty"`
}

func CommentAction(c *gin.Context) {

	getUserId, _ := c.Get("user_id")
	var userId int64
	if v, ok := getUserId.(int); ok {
		userId = int64(v)
	}

	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

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
		PostComment(c, userId, text, videoId)
	} else if actionType == "2" { //delete
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
		DeleteComment(c, videoId, commentId)
	}

}

func PostComment(c *gin.Context, userId int64, text string, videoId int64) {

	newComment := model.Comment{
		VideoId: videoId,
		UserId:  int64(userId),
		Content: text,
	}

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

	if err1 != nil || err2 != nil {
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
			ID:         int64(newComment.ID),
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: UserRespoonse{
				ID:            int64(getUser.ID),
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      false, // 需判断当前用户是否关注了视频作者，待其他servic完善后此处再完善
			},
		},
	})
}

func DeleteComment(c *gin.Context, videoId int64, commentId int64) {

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

	videoId := c.Query("video_id")
	commentList, err := service.GetCommentList(videoId)

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
		getUser, err := service.GetUser(commentList[i].UserId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseCommentList[i] = CommentResponse{
			ID:         int64(commentList[i].ID),
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"),
			User: UserRespoonse{
				ID:            int64(getUser.ID),
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      false, // 需判断当前用户是否关注了视频作者，待其他service完善后此处再完善
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
