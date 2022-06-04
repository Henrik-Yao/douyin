package service

import (
	"douyin/go/dao"
	"douyin/go/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

func RelationAction(relationreq *model.RelationRequest) (err error) {
	userId := relationreq.UserId
	touserId := relationreq.ToUserId
	actionType := relationreq.ActionType

	//数据准备
	relationAction := model.RelationAction{
		UserId:   userId,
		ToUserId: touserId,
	}
	dao.SqlSession.AutoMigrate(&model.RelationAction{}) //模型关联到数据库表videos
	if actionType == 1 {
		fmt.Println("执行了action_type == 1操作：")
		var relationExist = model.RelationAction{}
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
			//找不到数据
			//新建关注
			fmt.Println("关注不存在，新建关注")
			dao.SqlSession.Table("relations").Create(&relationAction)
			//user_login_infos的对应user_id的follow_count数量+1
			dao.SqlSession.Model(model.UserLoginInfo{}).Where("user_id=?", userId).Update("follow_count", gorm.Expr("follow_count+?", 1))
			//user_login_infos的对应to_user_id的follower_count数量+1
			dao.SqlSession.Model(model.UserLoginInfo{}).Where("user_id=?", touserId).Update("follower_count", gorm.Expr("follower_count+?", 1))
		} else {
			//找到数据
			fmt.Println("关注已存在")
		}

	}
	if actionType == 2 {
		fmt.Println("执行了action_type == 2操作：")
		var relationExist = model.RelationAction{}
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
			//找不到数据
			fmt.Println("关注不存在，无需操作")
			//dao.SqlSession.Table("relations").Create(&newrelation)
		} else {
			//找到数据
			fmt.Println("关注已存在，取消关注")
			//删除关注
			dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).Delete(&relationAction)
			//user_login_infos的对应user_id的follow_count数量-1
			dao.SqlSession.Model(model.UserLoginInfo{}).Where("user_id=?", userId).Update("follow_count", gorm.Expr("follow_count-?", 1))
			//user_login_infos的对应to_user_id的follower_count数量-1
			dao.SqlSession.Model(model.UserLoginInfo{}).Where("user_id=?", touserId).Update("follower_count", gorm.Expr("follower_count-?", 1))
		}
	}
	return nil
}
