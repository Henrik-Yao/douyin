package main

import (
	"douyin/go/dao"
	"douyin/go/model"
	"douyin/go/routes"
)

func main() {
	//连接数据库
	err := dao.InitMySql()
	if err != nil {
		panic(err)
	}
	//程序退出关闭数据库连接
	defer dao.Close()
	//绑定模型
	dao.SqlSession.AutoMigrate(&model.User{})
	//建数据库点赞记录表（用户id、视频id）
	dao.SqlSession.AutoMigrate(&model.FavoriteAction{})
	//注册路由
	r := routes.InitRouter()
	//启动端口为8080的项目
	r.Run(":8080")
}
