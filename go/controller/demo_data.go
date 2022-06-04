package controller

import "douyin/go/model"

var DemoVideos = []model.FeedVideo{
	{

		Author:  DemoUserinfo,
		PlayUrl: "http://" + "150.158.44.75" + ":" + "8080" + "/static/" + "bear.mp4",
	},
}
var DemoUserinfo = model.Author{

	Name:          "ava",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
