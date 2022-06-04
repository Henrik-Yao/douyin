package model

import (
	"github.com/jinzhu/gorm"
)

//用于数据库
type Relation struct {
	gorm.Model
	UserId   int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"`
}

//用于取数据，关注者/被关注者信息
type Follower struct {
	UserId        int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
