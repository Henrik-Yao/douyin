package service

import (
	"douyin/go/dao"
	"douyin/go/middleware"
	"douyin/go/model"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	MinPasswordLength = 8  //密码最小长度
)

type PostUserLoginFlow struct {
	username string
	password string

	data   *UserLoginResponse
	userid int64
	token  string
}

// PostUserLogin 注册用户并得到token和id
func PostUserLogin(username, password string) (*UserLoginResponse, error) {
	return NewPostUserLoginFlow(username, password).Do()
}

func NewPostUserLoginFlow(username, password string) *PostUserLoginFlow {
	//密码加密
	newPassword := dao.EncryptPassword([]byte(password))
	return &PostUserLoginFlow{username: username, password: newPassword}
}

func (q *PostUserLoginFlow) Do() (*UserLoginResponse, error) {
	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}

	//更新数据到数据库，也就是注册功能
	if err := q.updateData(); err != nil {
		return nil, err
	}

	//打包response
	if err := q.packResponse(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *PostUserLoginFlow) checkNum() error {
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

func (q *PostUserLoginFlow) updateData() error {

	//准备好userInfo,默认name为username,将加密密码写进数据库
	userLogin := model.UserLogin1{Username: q.username, Password: q.password}
	userinfo := model.UserInfo1{User: &userLogin, Name: q.username}

	//判断用户名是否已经存在
	userLoginDAO := dao.NewUserLoginDao()
	if userLoginDAO.IsUserExistByUsername(q.username) {
		return ErrorUserExit
	}

	//更新操作，由于userLogin属于userInfo，故更新userInfo即可，且由于传入的是指针，所以插入的数据内容也是清楚的
	userInfoDAO := dao.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userinfo)
	if err != nil {
		return err
	}

	//颁发token
	token, err := middleware.CreateToken(1, userLogin.Username)
	if err != nil {
		return err
	}
	q.token = token
	q.userid = userinfo.Id
	return nil
}

func (q *PostUserLoginFlow) packResponse() error {
	q.data = &UserLoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
