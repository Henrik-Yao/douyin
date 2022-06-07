package service

import (
	"douyin/src/dao"
	"douyin/src/model"
	"fmt"

	"github.com/jinzhu/gorm"
)

// Follower 用于取数据，关注者/被关注者信息
type Follower struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

//关注表
var followings = "followings"

//用户表

// IsFollowing 判断HostId是否关注GuestId
func IsFollowing(HostId uint, GuestId uint) bool {
	var relationExist = &model.Following{}
	if err := dao.SqlSession.Model(&model.Following{}).Where("host_id=? AND guest_id=?", HostId, GuestId).First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
		//关注不存在
		return false
	}
	//关注存在
	return true
}

// IncreaseFollowCount 增加HostId的关注数（Host_id 的 follow_count+1）
func IncreaseFollowCount(HostId uint) (err error) {
	dao.SqlSession.Model(&model.User{}).Where("id=?", HostId).Update("follow_count", gorm.Expr("follow_count+?", 1))
	return nil
}

// DecreaseFollowCount 增加HostId的关注数（Host_id 的 follow_count-1）
func DecreaseFollowCount(HostId uint) (err error) {
	dao.SqlSession.Model(&model.User{}).Where("id=?", HostId).Update("follow_count", gorm.Expr("follow_count-?", 1))
	return nil
}

// CreateFollowing 创建关注
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

// DeleteFollowing 删除关注
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

// FollowAction 关注操作
func FollowAction(HostId uint, GuestId uint, actionType uint) error {
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
			err := CreateFollowing(HostId, GuestId)
			if err != nil {
				return err
			}
			err = CreateFollower(GuestId, HostId)
			if err != nil {
				return err
			}
			//增加host_id的关注数
			err = IncreaseFollowCount(HostId)
			if err != nil {
				return err
			}
			//增加guest_id的粉丝数
			err = IncreaseFollowerCount(GuestId)
			if err != nil {
				return err
			}
		}
	}
	if actionType == 2 {
		//判断关注是否存在
		if IsFollowing(HostId, GuestId) {
			//关注存在
			fmt.Println("关注已存在,删除关注")
			//删除关注
			err := DeleteFollowing(HostId, GuestId)
			if err != nil {
				return err
			}
			err = DeleteFollower(GuestId, HostId)
			if err != nil {
				return err
			}
			//减少host_id的关注数
			err = DecreaseFollowCount(HostId)
			if err != nil {
				return err
			}
			//减少guest_id的粉丝数
			err = DecreaseFollowerCount(GuestId)
			if err != nil {
				return err
			}
		} else {
			//关注不存在
			fmt.Println("关注不存在")
		}
	}
	return nil
}

// FollowingList 获取关注表
func FollowingList(HostId uint) ([]Follower, error) {
	//1.followingList数据模型准备
	var followingList []Follower
	//var test []model.User

	//2.查HostId的关注表
	if err := dao.SqlSession.Model(&model.User{}).Joins("left join "+followings+" on "+users+".id = "+followings+".guest_id").
		Where(followings+".host_id=?", HostId).Scan(&followingList).Error; err != nil {
		return followingList, nil
	}
	fmt.Println(followingList)

	//3.修改查询结果中的is_follow属性
	for i, m := range followingList {
		if IsFollowing(m.Id, HostId) {
			//没有发生错误：找到
			fmt.Println("找到")
			followingList[i].IsFollow = true
		} else {
			//发生错误：没有找到
			fmt.Println("没找到")
			followingList[i].IsFollow = false
		}

	}

	return followingList, nil
}
