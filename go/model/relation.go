package model

import (
	"github.com/jinzhu/gorm"
)

type Relation struct {
	gorm.Model
	UserId   int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"`
}

type RelationAction struct {
	gorm.Model
	UserId   int64 `json:"user_id"`
	ToUserId int64 `json:"to_user_id"`
}

type RelationRequest struct {
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
	ToUserId   int64  `json:"to_user_id"`
	ActionType int32  `json:"action_type"`
}
