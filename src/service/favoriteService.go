package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

//查询某用户是否点赞某视频
func CheckFavorite(uid uint, vid uint) bool {
	var total int
	if err := dao.SqlSession.Table("favorites").
		Where("user_id = ? AND video_id = ? AND state = 1", uid, vid).Count(&total).
		Error; gorm.IsRecordNotFoundError(err) { //没有该条记录
		return false
	}
	if total == 0 {
		return false
	}
	return true
}

//点赞操作
//func FavoriteAction(favoritereq *controller.FavoriteRequest) (err error) {
	func FavoriteAction(userId uint,videoId uint,actionType uint) (err error) {
	
		//1-点赞
		if actionType == 1 {
			fmt.Println("执行了action_type == 1操作：")
			favoriteAction := model.Favorite{
				UserId:  userId,
				VideoId: videoId,
				State:   1,//1-已点赞
			}
			//var favorite_exist *model.Favorite//不对
			var favoriteExist = &model.Favorite{} //找不到时会返回错误
			//var favorite_exist []model.Favorite
			fmt.Println("执行了favorite_exist操作：")
			//如果没有记录-Create，如果有了记录-修改State
			result := dao.SqlSession.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteExist)
			fmt.Println("执行了查找操作：")
			fmt.Println(result.Error)
			if result.Error != nil { //不存在
				if err := dao.SqlSession.Table("favorites").Create(&favoriteAction).Error; err != nil { //创建记录
					fmt.Println("执行了创建操作：")
					fmt.Println(err)
					return err
				}
				dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			} else { //存在
				if favoriteExist.State == 0 {//state为0-video的favorite_count加1
					dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
					dao.SqlSession.Table("favorites").Where("video_id = ?", videoId).Update("state", 1)
				}
				//state为1-video的favorite_count不变
				return nil
			}
	
		} else { //2-取消点赞
			fmt.Println("执行了action_type == 2操作：")
			var favoriteCancel = &model.Favorite{}
			favoriteActionCancel := model.Favorite{
				UserId:  userId,
				VideoId: videoId,
				State:   0,//0-未点赞
			}
			if err := dao.SqlSession.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteCancel).Error; err != nil { //找不到这条记录，取消点赞失败，创建记录
				dao.SqlSession.Table("favorites").Create(&favoriteActionCancel)
				return err
			}
			fmt.Println(favoriteCancel)
			//存在
			if favoriteCancel.State == 1 {//state为1-video的favorite_count减1
				dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - 1"))
				//更新State
				dao.SqlSession.Table("favorites").Where("video_id = ?", videoId).Update("state", 0)
			}
			//state为0-video的favorite_count不变
			return nil
		}
		return nil
	}
	
	//获取点赞列表
	func FavoriteList(userId uint) ([]model.Video, error) {
	
		//查询当前id用户的所有点赞视频
		var favoriteList []model.Favorite
		videoList := make([]model.Video, 0)
		if err := dao.SqlSession.Table("favorites").Where("user_id=? AND state=?", userId,1).Find(&favoriteList).Error; err != nil { //找不到记录
			return videoList, nil
		}
		fmt.Println(favoriteList)//在favorites表中找到记录
		for _, m := range favoriteList {
			//dao.SqlSession.Table("videos").Where("id = ?", m.VideoId).Update("is_favorite", true) //is_favorite赋值为true
			var count int
			dao.SqlSession.Table("favorites").Where("video_id = ?", m.VideoId).Count(&count) //统计个数
			fmt.Println(count)
			dao.SqlSession.Table("videos").Where("id = ?", m.VideoId).Update("favorite_count", count)//更新favorite_count字段
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