package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"

	"github.com/jinzhu/gorm"
)

//获得视频列表
func FeedGet() ([]model.Video, error) {
	result := dao.SqlSession.Table("videos")
	fmt.Println("查询所有视频信息")
	fmt.Println(result.Error)
	var VideoList []model.Video
	VideoList = make([]model.Video, 0)
	err := dao.SqlSession.Table("videos").Find(&VideoList).Error
	return VideoList, err
}

// add comment_count
func AddCommentCount(videoId int64) error {

	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

// reduce comment_count
func ReduceCommentCount(videoId int64) error {

	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}
