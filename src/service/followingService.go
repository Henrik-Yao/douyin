package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"

	"github.com/jinzhu/gorm"
)

//用于取数据，关注者/被关注者信息
type Follower struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

//关注表
var followings string = "followings"

//用户表

// 判断HostId是否关注GuestId
func IsFollowing(HostId uint, GuestId uint) bool {
	var relationExist = &model.Following{}
	if err := dao.SqlSession.Model(&model.Following{}).Where("host_id=? AND guest_id=?", HostId, GuestId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
		//关注不存在
		return false
	}
	//关注存在
	return true
}

// 增加HostId的关注数（Host_id 的 follow_count+1）
func IncreaseFollowCount(HostId uint) (err error) {
	dao.SqlSession.Model(&model.User{}).Where("id=?", HostId).Update("follow_count", gorm.Expr("follow_count+?", 1))
	return nil
}

// 增加HostId的关注数（Host_id 的 follow_count-1）
func DecreaseFollowCount(HostId uint) (err error) {
	dao.SqlSession.Model(&model.User{}).Where("id=?", HostId).Update("follow_count", gorm.Expr("follow_count-?", 1))
	return nil
}

// 创建关注
func CreateFollowing(HostId uint, GuestId uint) (err error) {

	//1.Following数据模型准备
	newFollowing := model.Following{
		HostId:  HostId,
		GuestId: GuestId,
	}

	//2.模型关联到数据库表followings
	//dao.SqlSession.AutoMigrate(&model.Following{})

	//3.新建following
	dao.SqlSession.Model(&model.Following{}).Create(&newFollowing)
	return nil
}

// 删除关注
func DeleteFollowing(HostId uint, GuestId uint) (err error) {
	//1.Following数据模型准备
	newFollowing := model.Following{
		HostId:  HostId,
		GuestId: GuestId,
	}

	//2.模型关联到数据库表followings
	//dao.SqlSession.AutoMigrate(&model.Following{})

	//3.删除following
	dao.SqlSession.Model(&model.Following{}).Where("host_id=? AND guest_id=?", HostId, GuestId).Delete(&newFollowing)

	return nil
}

//关注操作
func FollowAction(HostId uint, GuestId uint, actionType uint) (err error) {
	//创建关注操作
	if actionType == 1 {
		//判断关注是否存在
		if IsFollowing(HostId, GuestId) {
			//关注存在
			fmt.Println("关注已存在")
		} else {
			//关注不存在
			fmt.Println("关注不存在，创建关注")
			//创建关注
			CreateFollowing(HostId, GuestId)
			CreateFollower(GuestId, HostId)
			//增加host_id的关注数
			IncreaseFollowCount(HostId)
			//增加guest_id的粉丝数
			IncreaseFollowerCount(GuestId)
		}
	}
	if actionType == 2 {
		//判断关注是否存在
		if IsFollowing(HostId, GuestId) {
			//关注存在
			fmt.Println("关注已存在,删除关注")
			//删除关注
			DeleteFollowing(HostId, GuestId)
			DeleteFollower(GuestId, HostId)
			//减少host_id的关注数
			DecreaseFollowCount(HostId)
			//减少guest_id的粉丝数
			DecreaseFollowerCount(GuestId)
		} else {
			//关注不存在
			fmt.Println("关注不存在")
		}
	}
	return nil
}

//获取关注表
func FollowingList(HostId uint) ([]Follower, error) {
	//1.followlist数据模型准备
	var followinglist []Follower
	//var test []model.User

	//2.查HostId的关注表
	if err := dao.SqlSession.Model(&model.User{}).Joins("left join "+followings+" on "+users+".id = "+followings+".guest_id").
		Where(followings+".host_id=?", HostId).Scan(&followinglist).Error; err != nil {
		return followinglist, nil
	}
	fmt.Println(followinglist)

	//3.修改查询结果中的is_follow属性
	for i, m := range followinglist {
		if IsFollowing(uint(m.Id), HostId) {
			//没有发生错误：找到
			fmt.Println("找到")
			followinglist[i].IsFollow = true
		} else {
			//发生错误：没有找到
			fmt.Println("没找到")
			followinglist[i].IsFollow = false
		}

	}

	return followinglist, nil
}
