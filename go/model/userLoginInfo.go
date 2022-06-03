package model

import "github.com/jinzhu/gorm"

type UserLoginInfo struct {
	gorm.Model
	UserId        int64  `json:"user_id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// UserLogin 用户登录表，和UserInfo属于一对一关系
type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:25;notnull"`
}
