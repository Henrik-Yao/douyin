package service

import (
	"douyin/go/dao"
	"douyin/go/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

//用户信息表名
var userdata string = "user_login_infos"

//关注操作
func RelationAction(userId int64, touserId int64, actionType int32) (err error) {

	//1.relation数据模型准备
	relation := model.Relation{
		UserId:   userId,
		ToUserId: touserId,
	}

	//2.模型关联到数据库表relations
	dao.SqlSession.AutoMigrate(&model.Relation{})

	//3.关注操作
	if actionType == 1 {
		fmt.Println("执行了action_type == 1操作：")

		var relationExist = model.Relation{}
		//3.1判断关注是否已经存在
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
			//关注不存在，新建关注
			fmt.Println("关注不存在，新建关注")
			dao.SqlSession.Table("relations").Create(&relation)
			//user_infos的对应user_id的follow_count数量+1
			dao.SqlSession.Table(userdata).Where("user_id=?", userId).Update("follow_count", gorm.Expr("follow_count+?", 1))
			//user_infos的对应to_user_id的follower_count数量+1
			dao.SqlSession.Table(userdata).Where("user_id=?", touserId).Update("follower_count", gorm.Expr("follower_count+?", 1))
		} else {
			//关注已存在，不做任何操作
			fmt.Println("关注已存在")
		}

	}
	//4.取消关注操作
	if actionType == 2 {
		fmt.Println("执行了action_type == 2操作：")
		var relationExist = model.Relation{}
		//4.1判断关注是否已经存在
		if err := dao.SqlSession.Table("relations").Where("user_id=? AND to_user_id=?", userId, touserId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
			//关注不存在，不做任何操作
			fmt.Println("关注不存在，无需操作")
		} else {
			//关注存在
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

//获取关注列表
func FollowList(userId int64) ([]model.Follower, error) {
	//1.followlist数据模型准备
	var followlist []model.Follower

	//2.查用户的关注表
	if err := dao.SqlSession.Table(userdata).Joins("left join relations on "+userdata+".user_id = relations.to_user_id").
		Where("relations.user_id=?", userId).Scan(&followlist).Error; err != nil {
		return followlist, nil
	}
	//3.修改查询结果中的is_follow属性
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

//获取粉丝列表
func FollowerList(userId int64) ([]model.Follower, error) {
	//1.followerlist数据模型准备
	var followerlist []model.Follower
	//2.查用户的粉丝表
	if err := dao.SqlSession.Table(userdata).Joins("left join relations on "+userdata+".user_id = relations.user_id").
		Where("relations.to_user_id=?", userId).Scan(&followerlist).Error; err != nil {
		return followerlist, nil
	}

	//3.修改查询结果中的is_follow属性
	for i, m := range followerlist {
		if err := dao.SqlSession.Table("relations").Where("to_user_id=? AND user_id=?", m.UserId, userId).Find(&model.Relation{}).Error; err != nil {
			//找到
			fmt.Println("没找到")
			followerlist[i].IsFollow = false
		} else {
			fmt.Println("找到")
			followerlist[i].IsFollow = true
		}

	}

	return followerlist, nil
}
