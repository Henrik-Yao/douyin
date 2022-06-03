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
