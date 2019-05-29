package common

import (
	"errors"
	"time"
	"yanfei_backend/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// 自定义错误
var (
	ErrTokenExpired = errors.New("Token is expired")
	ErrTokenInvalid = errors.New("Token is invalid")
)

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(viper.GetString("wechat.salt")),
	}
}

// ParseToken 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			return nil, ErrTokenInvalid
		}
	}
	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// CreateToken 生成令牌
func CreateToken(userID int64) (string, error) {
	j := NewJWT()
	claims := model.CustomClaims{
		userID,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix()),               // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 30*24*60*60), // 过期时间 一个月
			Issuer:    "LogicJake",                            //签名的发行者
		},
	}

	tokenNoSigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenNoSigned.SignedString(j.SigningKey)

	return token, err
}
