package common

import (
	"errors"
	"net/http"
	"regexp"
	"time"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var whiteList = []string{"/docs", "/wx/user/login"}
		var requestURL = c.Request.RequestURI

		for _, v := range whiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Next()
				return
			}
		}

		token := c.Request.Header.Get("token")

		if token == "" {
			c.JSON(http.StatusOK, controller.Message{
				Status: NoToken,
				Msg:    Errors[NoToken],
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				c.JSON(http.StatusOK, controller.Message{
					Status: TokenExpired,
					Msg:    Errors[TokenExpired],
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusServiceUnavailable, controller.Message{
				Status: TokenInvalid,
				Msg:    Errors[TokenInvalid],
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
		c.Next()
	}
}

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
