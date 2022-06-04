package model

// UserLogin1 用户登录表，和UserInfo1属于一对一关系
type UserLogin1 struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:25;notnull"`
}
