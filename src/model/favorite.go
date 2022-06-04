package model

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	UserId  int32 `json:"user_id"`
	VideoId int32 `json:"video_id"`
	State   int32
}
