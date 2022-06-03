package controller


import (
	"douyin/go/middleware"
	"douyin/go/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
    "douyin/go/service"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type FavoriteListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

//点赞视频方法
func Favorite(c *gin.Context) {
	//验证token
	//参数绑定
    var favoritereq model.FavoriteRequest
	c.BindJSON(&favoritereq)
    token := favoritereq.Token
	//token验证：
	tokenStruck, ok := middleware.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "token不正确"})
		c.Abort() //阻止执行
		return
	}
	//token超时
	if time.Now().Unix() > tokenStruck.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "token过期"})
		c.Abort() //阻止执行
		return
	}
	fmt.Println("token通过验证")
   
	//函数调用及响应
	err := service.FavoriteAction(&favoritereq)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				StatusMsg:  "操作成功！",
			})
		}
}



//获取列表方法
func FavoriteList(c *gin.Context) {
	//鉴权token
	user_id := c.Query("user_id")
	token := c.Query("token")

	tokenStruck, ok := middleware.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 403, "msg": "token不正确"})
		c.Abort() //阻止执行
		return
	}
	//token超时
	if time.Now().Unix() > tokenStruck.ExpiresAt {
		c.JSON(http.StatusOK, gin.H{"code": 402, "msg": "token过期"})
		c.Abort() //阻止执行
		return
	}
	fmt.Println("token通过验证")

	
	//函数调用及响应
	video_list, err := service.FavoriteList(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, FavoriteListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "查找列表失败！",
				},
				VideoList: nil,
			})
		} else {
			c.JSON(http.StatusOK, FavoriteListResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "已找到列表！",
				},
				VideoList: video_list,
			})
		}
 }