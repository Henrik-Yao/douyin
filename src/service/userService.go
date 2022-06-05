package service

import (
	"douyin/src/common"
	"douyin/src/dao"
	"douyin/src/middleware"
	"douyin/src/model"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	MinPasswordLength = 8  //密码最小长度
)

func CreateUser(user *model.User) (err error) {
	if err = dao.SqlSession.Create(user).Error; err != nil {
		return err
	}
	return
}

type UserLoginResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type QueryUserLoginFlow struct {
	username string
	password string

	data   *UserLoginResponse
	userid uint
	token  string
}

// UserLogin 登录功能，查询用户是否存在，并返回token和id
func UserLogin(username, password string) (*UserLoginResponse, error) {
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
		return common.ErrorUserNameNull
	}
	if len(q.username) > MaxUsernameLength {
		return common.ErrorUserNameExtend
	}
	if q.password == "" {
		return common.ErrorPasswordNull
	}
	if len(q.password) > MaxPasswordLength || len(q.password) < MinPasswordLength {
		return common.ErrorPasswordLength
	}
	return nil
}

func (q *QueryUserLoginFlow) prepareData() error {
	userLoginDAO := model.NewUserLoginDao()
	var login model.User
	//准备好userid
	err := userLoginDAO.QueryUserLogin(q.username, q.password, &login)
	if err != nil {
		return err
	}
	q.userid = login.Model.ID

	//准备颁发token
	token, err := middleware.CreateToken(login.Model.ID, login.Name)
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

type PostUserLoginFlow struct {
	username string
	password string

	data   *UserLoginResponse
	userid uint
	token  string
}

//UserRegister 注册用户并得到token和id
func UserRegister(username, password string) (*UserLoginResponse, error) {
	return NewPostUserLoginFlow(username, password).Do()
}

func NewPostUserLoginFlow(username, password string) *PostUserLoginFlow {
	//密码加密
	newPassword := model.EncryptPassword([]byte(password))
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
		return common.ErrorUserNameNull
	}
	if len(q.username) > MaxUsernameLength {
		return common.ErrorUserNameExtend
	}
	if q.password == "" {
		return common.ErrorPasswordNull
	}
	if len(q.password) > MaxPasswordLength || len(q.password) < MinPasswordLength {
		return common.ErrorPasswordLength
	}
	return nil
}

func (q *PostUserLoginFlow) updateData() error {

	user := model.User{Name: q.username, Password: q.password}
	//判断用户名是否已经存在
	userLoginDAO := model.NewUserLoginDao()
	if userLoginDAO.IsUserExistByUsername(q.username) {
		return common.ErrorUserExit
	}

	//更新操作，由于传入的是指针，所以插入的数据内容也是清楚的
	userInfoDAO := model.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&user)
	if err != nil {
		return err
	}

	//颁发token
	token, err := middleware.CreateToken(user.ID, user.Name)
	if err != nil {
		return err
	}
	q.token = token
	q.userid = user.Model.ID
	return nil

}

func (q *PostUserLoginFlow) packResponse() error {
	q.data = &UserLoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
