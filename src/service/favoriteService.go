package service

import (
	"douyin/src/dao"
	"github.com/jinzhu/gorm"
)

//查询某用户是否点赞某视频
func CheckFavorite(uid int, vid int) bool {
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
