package service

import (
	"douyin/go/dao"
	"douyin/go/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

//点赞操作
func FavoriteAction(favoritereq *model.FavoriteRequest) (err error) {
	//参数获取
	userId := favoritereq.UserId
	videoId := favoritereq.VideoId
	actionType := favoritereq.ActionType
	//不能重复点赞
	//1-点赞
	if actionType == 1 {
		fmt.Println("执行了action_type == 1操作：")
		favoriteAction := model.FavoriteAction{
			UserId:  userId,
			VideoId: videoId,
		}
		//var favorite_exist *model.FavoriteAction//不对
		var favoriteExist = &model.FavoriteAction{} //找不到时会返回错误
		//var favorite_exist []model.FavoriteAction
		fmt.Println("执行了favorite_exist操作：")
		result := dao.SqlSession.Table("favorite_actions").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteExist)
		fmt.Println("执行了查找操作：")
		fmt.Println(result.Error)
		if result.Error != nil { //不存在
			if err := dao.SqlSession.Table("favorite_actions").Create(&favoriteAction).Error; err != nil { //创建记录
				fmt.Println("执行了创建操作：")
				fmt.Println(err)
				return err
			}
			dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
		} else {
			return nil
		}

	} else { //2-取消点赞
		fmt.Println("执行了action_type == 2操作：")
		var favoriteCancel = &model.FavoriteAction{}
		if err := dao.SqlSession.Table("favorite_actions").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteCancel).Error; err != nil { //找不到这条记录，取消点赞失败
			return err
		}
		fmt.Println(favoriteCancel)
		var favoriteAction = &model.FavoriteAction{}
		//var favorite_action *model.FavoriteAction//记录存在，删除记录
		if err := dao.SqlSession.Table("favorite_actions").Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&favoriteAction).Error; err != nil {
			return err
		}
		dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - 1"))

		return nil
	}
	return nil
}

//获取点赞列表
func FavoriteList(userId string) ([]model.Video, error) {

	//查询当前id用户的所有点赞视频
	var favoriteList []model.FavoriteAction
	videoList := make([]model.Video, 0)
	if err := dao.SqlSession.Table("favorite_actions").Where("user_id=?", userId).Find(&favoriteList).Error; err != nil { //找不到记录
		return videoList, nil
	} //user_id-string类型
	fmt.Println(favoriteList)
	for _, m := range favoriteList { //给video的is_favorite、favorite_count、is_comment、comment_count字段赋值
		dao.SqlSession.Table("videos").Where("id = ?", m.VideoId).Update("is_favorite", true) //is_favorite赋值为true
		var count int
		dao.SqlSession.Table("favorite_actions").Where("video_id = ?", m.VideoId).Count(&count) //统计个数
		fmt.Println(count)
		dao.SqlSession.Table("videos").Where("id = ?", m.VideoId).Update("favorite_count", count)
		//fmt.Println(m.VideoId)
		var video = model.Video{}
		if err := dao.SqlSession.Table("videos").Where("id=?", m.VideoId).Find(&video).Error; err != nil {
			return nil, err
		}
		fmt.Println(video)
		videoList = append(videoList, video)
		fmt.Println(videoList)
	}
	return videoList, nil
}
