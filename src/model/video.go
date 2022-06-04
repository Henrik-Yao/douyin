package model

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model
	AuthorId      int32  `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	Title         string `json:"title"`
}
