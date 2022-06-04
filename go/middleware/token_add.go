package middleware

import (
	"douyin/go/dao"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	return token, claims, err
}

// JWTAuth 用于验证token，并返回token对应的userid
func JWTAuth(token string) (int64, error) {
	if token == "" {
		return 0, errors.New("token为空")
	}
	_, claim, err := ParseToken(token)
	if err != nil {
		return 0, errors.New("token过期")
	}
	//最后验证这个user是否真的存在
	if !dao.NewUserInfoDAO().IsUserExistById(claim.UserId) {
		return 0, errors.New("user不存在")
	}
	return claim.UserId, nil
}
