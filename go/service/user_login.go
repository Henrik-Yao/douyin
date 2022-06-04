package service

import (
	"douyin/go/dao"
	"douyin/go/middleware"
	"douyin/go/model"
)

type UserLoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type QueryUserLoginFlow struct {
	username string
	password string

	data   *UserLoginResponse
	userid int64
	token  string
}

// QueryUserLogin 查询用户是否存在，并返回token和id
func QueryUserLogin(username, password string) (*UserLoginResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do()
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

func (q *QueryUserLoginFlow) Do() (*UserLoginResponse, error) {
	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	//准备好数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//打包最终数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *QueryUserLoginFlow) checkNum() error {
	if q.username == "" {
		return ErrorUserNameNull
	}
	if len(q.username) > MaxUsernameLength {
		return ErrorUserNameExtend
	}
	if q.password == "" {
		return ErrorPasswordNull
	}
	if len(q.password) > MaxPasswordLength || len(q.password) < MinPasswordLength {
		return ErrorPasswordLength
	}
	return nil
}

func (q *QueryUserLoginFlow) prepareData() error {
	userLoginDAO := dao.NewUserLoginDao()
	var login model.UserLogin
	//准备好userid
	err := userLoginDAO.QueryUserLogin(q.username, q.password, &login)
	if err != nil {
		return err
	}
	q.userid = login.Id

	//准备颁发token
	token, err := middleware.CreateToken(login.Id, login.Username)
	if err != nil {
		return err
	}
	q.token = token
	return nil
}

func (q *QueryUserLoginFlow) packData() error {
	q.data = &UserLoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
