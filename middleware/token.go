package middleware

import (
	"regexp"
	"hackthoon/common"

	"github.com/gin-gonic/gin"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var whiteList = []string{"/docs", "/wx/user/login", "/ping"}
		var requestURL = c.Request.RequestURI

		for _, v := range whiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Next()
				return
			}
		}

		token := c.Request.Header.Get("token")

		if common.FuncHandler(c, token != "", true, common.NoToken) {
			c.Abort()
			return
		}

		j := common.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)

		if common.FuncHandler(c, err != common.ErrTokenExpired, true, common.TokenExpired) {
			c.Abort()
			return
		}

		if common.FuncHandler(c, err != common.ErrTokenInvalid, true, common.TokenInvalid) {
			c.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
		c.Next()
	}
}
