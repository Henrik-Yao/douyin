package dao

import (
	"douyin/go/model"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once //单例
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

// QueryUserInfoById 根据用户ID查询用户信息
func (u *UserInfoDAO) QueryUserInfoById(userId int64, userinfo *model.UserInfo1) error {
	if userinfo == nil {
		return ErrorNullPointer
	}
	SqlSession.Where("id=?", userId).First(userinfo)
	//id为零值，说明sql执行失败
	if userinfo.Id == 0 {
		return ErrorUserNotExit
	}
	return nil
}

// AddUserInfo 添加用户信息
func (u *UserInfoDAO) AddUserInfo(userinfo *model.UserInfo1) error {
	if userinfo == nil {
		return ErrorNullPointer
	}
	return SqlSession.Create(userinfo).Error
}

// IsUserExistById 根据用户ID判断用户是否存在
func (u *UserInfoDAO) IsUserExistById(id int64) bool {
	var userinfo model.UserInfo1
	if err := SqlSession.Where("id=?", id).Select("id").First(&userinfo).Error; err != nil {
		log.Println(err)
	}
	if userinfo.Id == 0 {
		return false
	}
	return true
}

// CancelUserFollow 根据userId和被关注者userToId，实现取关
func (u *UserInfoDAO) CancelUserFollow(userId, userToId int64) error {
	return SqlSession.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", userToId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_relations` WHERE user_info_id=? AND follow_id=?", userId, userToId).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetFollowListByUserId 根据userId获得关注列表
func (u *UserInfoDAO) GetFollowListByUserId(userId int64, userList *[]*model.UserInfo1) error {
	if userList == nil {
		return ErrorNullPointer
	}
	var err error
	if err = SqlSession.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.user_info_id = ? AND r.follow_id = u.id", userId).Scan(userList).Error; err != nil {
		return err
	}
	if (*userList)[0].Id == 0 {
		return ErrEmptyUserList
	}
	return nil
}

// GetFollowerListByUserId 根据userId得到粉丝列表
func (u *UserInfoDAO) GetFollowerListByUserId(userId int64, userList *[]*model.UserInfo1) error {
	if userList == nil {
		return ErrorNullPointer
	}
	var err error
	if err = SqlSession.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.follow_id = ? AND r.user_info_id = u.id", userId).Scan(userList).Error; err != nil {
		return err
	}
	//if len(*userList) == 0 || (*userList)[0].Id == 0 {
	//	return ErrEmptyUserList
	//}
	return nil
}
