package model

import "github.com/jinzhu/gorm"

type Comment struct { // 评论
	gorm.Model
	VideoId int64  `json:"video_id,omitempty"`
	UserId  int64  `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}
