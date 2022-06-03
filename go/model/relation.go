package model

import (
	"github.com/jinzhu/gorm"
)

type Relation struct {
	gorm.Model
	UserId   int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"`
}
