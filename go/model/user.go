package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName string `json:"username"`
	Password string `json:"password"`
}
