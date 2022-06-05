package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"
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
