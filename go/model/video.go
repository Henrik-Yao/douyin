package model

import "github.com/jinzhu/gorm"

type Video struct { //视频信息
	VideoId       int64     `json:"video_id,omitempty"`
	Author        LoginUser `json:"author"`
	PlayUrl       string    `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
	Title         string    `json:"title"`
}

type VideoNoAuthor struct {
	gorm.Model
	UserId        int64  `json:"user_id"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title"`
}
