package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"time"
)

// GetCommentList 获取指定videoId的评论表
func GetCommentList(videoId uint) ([]model.Comment, error) {
	var commentList []model.Comment
	if err := dao.SqlSession.Table("comments").Where("video_id=?", videoId).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// PostComment 发布评论
func PostComment(comment model.Comment) error {
	if err := dao.SqlSession.Table("comments").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteComment 删除指定commentId的评论
func DeleteComment(commentId uint) error {
	if err := dao.SqlSession.Table("comments").Where("id = ?", commentId).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
