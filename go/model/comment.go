package model

import "github.com/jinzhu/gorm"

type Comment struct { // 评论
	gorm.Model
	Id         int64  `json:"id,omitempty"`
	UserId     int64  `json:"user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoId    int64  `json:"video_id,omitempty"`
	IsDeleted  bool   `json:"is_deleted,omitempty"`
}
