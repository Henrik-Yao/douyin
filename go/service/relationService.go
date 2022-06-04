package service

import (
	"douyin/go/dao"
	"douyin/go/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

var userdata string = "user_login_infos"

func RelationAction(userId int64, touserId int64, actionType int32) (err error) {

	//数据准备
	relation := model.Relation{
		UserId:   userId,
		ToUserId: touserId,
	}
	dao.SqlSession.AutoMigrate(&model.Relation{}) //模型关联到数据库表videos
	if actionType == 1 {
		fmt.Println("执行了action_type == 1操作：")
		var relationExist = model.Relation{}
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
			//找不到数据
			//新建关注
			fmt.Println("关注不存在，新建关注")
			dao.SqlSession.Table("relations").Create(&relation)
			//user_infos的对应user_id的follow_count数量+1
			dao.SqlSession.Table(userdata).Where("user_id=?", userId).Update("follow_count", gorm.Expr("follow_count+?", 1))
			//user_infos的对应to_user_id的follower_count数量+1
			dao.SqlSession.Table(userdata).Where("user_id=?", touserId).Update("follower_count", gorm.Expr("follower_count+?", 1))
		} else {
			//找到数据
			fmt.Println("关注已存在")
		}

	}
	if actionType == 2 {
		fmt.Println("执行了action_type == 2操作：")
		var relationExist = model.Relation{}
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
			//找不到数据
			fmt.Println("关注不存在，无需操作")
			//dao.SqlSession.Table("relations").Create(&newrelation)
		} else {
			//找到数据
			fmt.Println("关注已存在，取消关注")
			//删除关注
			dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).Delete(&relation)
			//user_infos的对应user_id的follow_count数量-1
			dao.SqlSession.Table(userdata).Where("user_id=?", userId).Update("follow_count", gorm.Expr("follow_count-?", 1))
			//user_infos的对应to_user_id的follower_count数量-1
			dao.SqlSession.Table(userdata).Where("user_id=?", touserId).Update("follower_count", gorm.Expr("follower_count-?", 1))
		}
	}
	return nil
}

func FollowList(userId int64) ([]model.Follower, error) {
	//2.先查relation表
	var followlist []model.Follower
	if err := dao.SqlSession.Table(userdata).Joins("left join relations on "+userdata+".user_id = relations.to_user_id").
		Where("relations.user_id=?", userId).Scan(&followlist).Error; err != nil {
		return followlist, nil
	}

	for i, m := range followlist {
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", m.UserId, userId).Find(&model.Relation{}).Error; err != nil {
			//发生错误：没有找到
			fmt.Println("没找到")
			followlist[i].IsFollow = false
		} else {
			//没有发生错误：找到
			fmt.Println("找到")
			followlist[i].IsFollow = true
		}

	}

	return followlist, nil
}

func FollowerList(userId int64) ([]model.Follower, error) {
	//2.先查relation表
	var followlist []model.Follower
	if err := dao.SqlSession.Table(userdata).Joins("left join relations on "+userdata+".user_id = relations.user_id").
		Where("relations.to_user_id=?", userId).Scan(&followlist).Error; err != nil {
		return followlist, nil
	}

	for i, m := range followlist {
		if err := dao.SqlSession.Table("relations").Where("to_user_id=? AND user_id=?", m.UserId, userId).Find(&model.Relation{}).Error; err != nil {
			//找到
			fmt.Println("没找到")
			followlist[i].IsFollow = false
		} else {
			fmt.Println("找到")
			followlist[i].IsFollow = true
		}

	}

	return followlist, nil
}
