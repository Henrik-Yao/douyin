package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"time"
)

func GetCommentList(videoId uint) ([]model.Comment, error) {
	var commentList []model.Comment
	if err := dao.SqlSession.Table("comment").Where("video_id=? and deleted_at is null", videoId).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

func PostComment(comment model.Comment) error {
	if err := dao.SqlSession.Table("comments").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentId uint) error {
	if err := dao.SqlSession.Table("comments").Where("id = ?", commentId).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
