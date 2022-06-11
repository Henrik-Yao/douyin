package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name           string `json:"name"`
	Password       string `json:"password"`
	FollowCount    uint   `json:"follow_count"`
	FollowerCount  uint   `json:"follower_count"`
	TotalFavorited uint   `json:"total_favorited"`
	FavoriteCount  uint   `json:"favorite_count"`
}
