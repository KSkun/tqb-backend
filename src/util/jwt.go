package util

import (
	"github.com/KSkun/tqb-backend/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	jwtExpiresDuration = time.Hour * 5 // 5 h
)

// JWTClaims 使用的JWT结构，JWT的修改请直接修改结构中的字段
type JWTClaims struct {
	jwt.StandardClaims
	User string `json:"user"`
}

// GenerateJWTToken 根据键值对生成jwt token
func GenerateJWTToken(claims JWTClaims) (string, int64, error) {
	claims.ExpiresAt = time.Now().Add(jwtExpiresDuration).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tSigned, err := t.SignedString([]byte(config.C.JWT.Secret))
	if err != nil {
		return tSigned, 0, err
	}
	return tSigned, claims.ExpiresAt, nil
}
