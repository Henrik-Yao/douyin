package controller

import (
	"douyin/src/common"
	"douyin/src/model"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateUser(c *gin.Context) {
	// 定义一个User变量
	var user model.User
	// 将调用后端的request请求中的body数据根据json格式解析到User结构变量中
	c.BindJSON(&user)
	// 将被转换的user变量传给service层的CreateUser方法，进行User的新建
	err := service.CreateUser(&user)
	// 判断是否异常，无异常则返回包含200和更新数据的信息
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"tip": "测试失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"tip": "测试成功",
		})
	}
}

//这里和commentController中的UserResponse重名，所以加个1避免重名

type UserResponse1 struct {
	common.Response
	User *model.User `json:"user"`
}

func UserInfo(c *gin.Context) {
	p := NewProxyUserInfo(c)
	//根据user_id查询
	rawId := c.Query("user_id")
	err := p.DoQueryUserInfoByUserId(rawId)
	//未发生错误，则就不用再使用token字段了
	if err == nil {
		return
	}

}

type ProxyUserInfo struct {
	c *gin.Context
}

func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}

func (p *ProxyUserInfo) DoQueryUserInfoByUserId(rawId string) error {
	userId, err := strconv.ParseInt(rawId, 10, 64)
	if err != nil {
		return err
	}
	//由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
	userinfoDAO := model.NewUserInfoDAO()

	var userInfo model.User
	err = userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse1{
		Response: common.Response{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) UserInfoOk(user *model.User) {
	p.c.JSON(http.StatusOK, UserResponse1{
		Response: common.Response{StatusCode: 0},
		User:     user,
	})
}

type UserLoginResponse struct {
	common.Response
	*service.UserLoginResponse
}

func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := service.UserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:          common.Response{StatusCode: 0},
		UserLoginResponse: userLoginResponse,
	})
}

type UserRegisterResponse struct {
	common.Response
	*service.UserLoginResponse
}

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	registerResponse, err := service.UserRegister(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response:          common.Response{StatusCode: 0},
		UserLoginResponse: registerResponse,
	})
}
