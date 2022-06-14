package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"github.com/jinzhu/gorm"
)

// CheckFavorite 查询某用户是否点赞某视频
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

// AddTotalFavorited 增加total_favorited
func AddTotalFavorited(HostId uint) error {
	if err := dao.SqlSession.Model(&model.User{}).
		Where("id=?", HostId).
		Update("total_favorited", gorm.Expr("total_favorited+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// ReduceTotalFavorited 减少total_favorited
func ReduceTotalFavorited(HostId uint) error {
	if err := dao.SqlSession.Model(&model.User{}).
		Where("id=?", HostId).
		Update("total_favorited", gorm.Expr("total_favorited-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// AddFavoriteCount 增加favorite_count
func AddFavoriteCount(HostId uint) error {
	if err := dao.SqlSession.Model(&model.User{}).
		Where("id=?", HostId).
		Update("favorite_count", gorm.Expr("favorite_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// ReduceFavoriteCount 减少favorite_count
func ReduceFavoriteCount(HostId uint) error {
	if err := dao.SqlSession.Model(&model.User{}).
		Where("id=?", HostId).
		Update("favorite_count", gorm.Expr("favorite_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// FavoriteAction 点赞操作
func FavoriteAction(userId uint, videoId uint, actionType uint) (err error) {

	//1-点赞
	if actionType == 1 {
		favoriteAction := model.Favorite{
			UserId:  userId,
			VideoId: videoId,
			State:   1, //1-已点赞
		}
		var favoriteExist = &model.Favorite{} //找不到时会返回错误
		//如果没有记录-Create，如果有了记录-修改State
		result := dao.SqlSession.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteExist)
		if result.Error != nil { //不存在
			if err := dao.SqlSession.Table("favorites").Create(&favoriteAction).Error; err != nil { //创建记录
				return err
			}
			dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
			//userId的favorite_count增加
			if err := AddFavoriteCount(userId); err != nil {
				return err
			}
			//videoId对应的userId的total_favorite增加
			GuestId, err := GetVideoAuthor(videoId)
			if err != nil {
				return err
			}
			if err := AddTotalFavorited(GuestId); err != nil {
				return err
			}
		} else { //存在
			if favoriteExist.State == 0 { //state为0-video的favorite_count加1
				dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
				dao.SqlSession.Table("favorites").Where("video_id = ?", videoId).Update("state", 1)
				//userId的favorite_count增加
				if err := AddFavoriteCount(userId); err != nil {
					return err
				}
				//videoId对应的userId的total_favorite增加
				GuestId, err := GetVideoAuthor(videoId)
				if err != nil {
					return err
				}
				if err := AddTotalFavorited(GuestId); err != nil {
					return err
				}
			}
			//state为1-video的favorite_count不变
			return nil
		}

	} else { //2-取消点赞
		var favoriteCancel = &model.Favorite{}
		favoriteActionCancel := model.Favorite{
			UserId:  userId,
			VideoId: videoId,
			State:   0, //0-未点赞
		}
		if err := dao.SqlSession.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteCancel).Error; err != nil { //找不到这条记录，取消点赞失败，创建记录
			dao.SqlSession.Table("favorites").Create(&favoriteActionCancel)
			//userId的favorite_count增加
			if err := ReduceFavoriteCount(userId); err != nil {
				return err
			}
			//videoId对应的userId的total_favorite增加
			GuestId, err := GetVideoAuthor(videoId)
			if err != nil {
				return err
			}
			if err := ReduceTotalFavorited(GuestId); err != nil {
				return err
			}
			return err
		}
		//存在
		if favoriteCancel.State == 1 { //state为1-video的favorite_count减1
			dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			//更新State
			dao.SqlSession.Table("favorites").Where("video_id = ?", videoId).Update("state", 0)
			if err := ReduceFavoriteCount(userId); err != nil {
				return err
			}
			//videoId对应的userId的total_favorite增加
			GuestId, err := GetVideoAuthor(videoId)
			if err != nil {
				return err
			}
			if err := ReduceTotalFavorited(GuestId); err != nil {
				return err
			}
			return err
		}
		//state为0-video的favorite_count不变
		return nil
	}
	return nil
}

// FavoriteList 获取点赞列表
func FavoriteList(userId uint) ([]model.Video, error) {

	//查询当前id用户的所有点赞视频
	var favoriteList []model.Favorite
	videoList := make([]model.Video, 0)
	if err := dao.SqlSession.Table("favorites").Where("user_id=? AND state=?", userId, 1).Find(&favoriteList).Error; err != nil { //找不到记录
		return videoList, nil
	}
	for _, m := range favoriteList {

		var video = model.Video{}
		if err := dao.SqlSession.Table("videos").Where("id=?", m.VideoId).Find(&video).Error; err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}
