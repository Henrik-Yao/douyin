package model

import "github.com/jinzhu/gorm"

type following struct {
	gorm.Model
	HostId  int32
	GuestId int32
}

type followers struct {
	gorm.Model
	HostId  int32
	GuestId int32
}
