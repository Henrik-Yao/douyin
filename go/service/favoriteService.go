package service

import (
	"douyin/go/dao"
	"douyin/go/model"
	"fmt"
	"github.com/jinzhu/gorm"
)


//点赞操作
func FavoriteAction(favoritereq *model.FavoriteRequest) (err error) {
	//参数获取
	user_id := favoritereq.UserId
	video_id := favoritereq.VideoId
	action_type := favoritereq.ActionType
	//不能重复点赞
	//1-点赞
	if action_type == 1{
		fmt.Println("执行了action_type == 1操作：")
		favorite_action := model.FavoriteAction{
			UserId:        user_id,
			VideoId:       video_id,
		}
		//var favorite_exist *model.FavoriteAction//不对
		var favorite_exist = &model.FavoriteAction{}//找不到时会返回错误
		//var favorite_exist []model.FavoriteAction
		fmt.Println("执行了favorite_exist操作：")
		result := dao.SqlSession.Table("favorite_actions").Where("user_id = ? AND video_id = ?", user_id, video_id).First(&favorite_exist)
		fmt.Println("执行了查找操作：")
		fmt.Println(result.Error)
		if result.Error != nil{//不存在
			if err := dao.SqlSession.Table("favorite_actions").Create(&favorite_action).Error; err != nil{//创建记录
				fmt.Println("执行了创建操作：")
				fmt.Println(err)
				return err 
			}
			dao.SqlSession.Table("videos").Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count + 1"))
		}else{
			return nil
		}

	
	}else{//2-取消点赞
		fmt.Println("执行了action_type == 2操作：")
		var favorite_cancel = &model.FavoriteAction{}
		if err :=dao.SqlSession.Table("favorite_actions").Where("user_id = ? AND video_id = ?", user_id, video_id).First(&favorite_cancel).Error; err != nil{//找不到这条记录，取消点赞失败
			return err
		}
		fmt.Println(favorite_cancel)
		var favorite_action = &model.FavoriteAction{}
		//var favorite_action *model.FavoriteAction//记录存在，删除记录
		if err := dao.SqlSession.Table("favorite_actions").Where("user_id = ? AND video_id = ?", user_id, video_id).Delete(&favorite_action).Error; err != nil{
			return err 
		}
		dao.SqlSession.Table("videos").Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count - 1"))

		return nil
	}
	return nil
}

//获取点赞列表
func FavoriteList(user_id string) ([]model.Video,  error) {
	
	//查询当前id用户的所有点赞视频
	var favorite_list []model.FavoriteAction
	video_list := make([]model.Video, 0)
	if err := dao.SqlSession.Table("favorite_actions").Where("user_id=?", user_id).Find(&favorite_list).Error; err != nil{//找不到记录
		return video_list, nil
	}//user_id-string类型
	fmt.Println(favorite_list)
	for _, m := range favorite_list {//给video的is_favorite、favorite_count、is_comment、comment_count字段赋值
		dao.SqlSession.Table("videos").Where("id = ?", m.VideoId).Update("is_favorite",true)//is_favorite赋值为true
		var count  int
		dao.SqlSession.Table("favorite_actions").Where("video_id = ?", m.VideoId).Count(&count)//统计个数
		fmt.Println(count)
		dao.SqlSession.Table("videos").Where("id = ?", m.VideoId).Update("favorite_count", count)
		//fmt.Println(m.VideoId)
		var video = model.Video{}
		if err := dao.SqlSession.Table("videos").Where("id=?", m.VideoId).Find(&video).Error; err != nil{
			return nil, err
		}
		fmt.Println(video)
		video_list = append(video_list, video)
		fmt.Println(video_list)
	}
		return video_list, nil
}