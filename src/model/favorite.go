package model

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	UserId  uint `json:"user_id"`
	VideoId uint `json:"video_id"`
	State   uint
}
