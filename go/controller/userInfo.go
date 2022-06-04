package controller

import (
	"douyin/go/dao"
	"douyin/go/middleware"
	"douyin/go/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserResponse struct {
	model.Response
	User *model.UserInfo1 `json:"user"`
}

func UserInfoHandler(c *gin.Context) {
	p := NewProxyUserInfo(c)
	//根据user_id查询
	rawId := c.Query("user_id")
	err := p.DoQueryUserInfoByUserId(rawId)
	//未发生错误，则就不用再使用token字段了
	if err == nil {
		return
	}

	//根据token查询
	token := c.Query("token")
	err = p.DoQueryUserInfoByToken(token)
	if err != nil {
		p.UserInfoError(err.Error())
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
	userinfoDAO := dao.NewUserInfoDAO()

	var userInfo model.UserInfo1
	err = userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *ProxyUserInfo) DoQueryUserInfoByToken(token string) error {
	userId, err := middleware.JWTAuth(token)
	if err != nil {
		return err
	}
	//由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
	userinfoDAO := dao.NewUserInfoDAO()

	var userInfo model.UserInfo1
	err = userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		Response: model.Response{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) UserInfoOk(user *model.UserInfo1) {
	p.c.JSON(http.StatusOK, UserResponse{
		Response: model.Response{StatusCode: 0},
		User:     user,
	})
}
