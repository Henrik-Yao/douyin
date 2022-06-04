package model

import "github.com/jinzhu/gorm"

type Following struct {
	gorm.Model
	HostId  int32
	GuestId int32
}

type Followers struct {
	gorm.Model
	HostId  int32
	GuestId int32
}
