package model

import "github.com/jinzhu/gorm"

type Following struct {
	gorm.Model
	HostId  uint
	GuestId uint
}

type Followers struct {
	gorm.Model
	HostId  uint
	GuestId uint
}
