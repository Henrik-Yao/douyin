package model

type UserInfo1 struct { // 视频发布者信息
	Id            int64       `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	FollowCount   int64       `json:"follow_count,omitempty"`
	FollowerCount int64       `json:"follower_count,omitempty"`
	IsFollow      bool        `json:"is_follow,omitempty"`
	User          *UserLogin1 `json:"-"` //用户与密码
}
