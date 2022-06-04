package model

import (
	"github.com/jinzhu/gorm"
)

type Relation struct {
	gorm.Model
	UserId   int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"`
}

type RelationRequest struct {
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
	ToUserId   int64  `json:"to_user_id"`
	ActionType int32  `json:"action_type"`
}

type Follower struct {
	UserId        int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
