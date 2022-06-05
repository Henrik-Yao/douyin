package model

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model
	AuthorId      uint   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount uint   `json:"favorite_count"`
	CommentCount  uint   `json:"comment_count"`
	Title         string `json:"title"`
}
