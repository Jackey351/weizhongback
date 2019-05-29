package middleware

import (
	"regexp"
	"strings"
	"yanfei_backend/common"
	"yanfei_backend/model"
	"yanfei_backend/storage"

	"github.com/gin-gonic/gin"
)

// Certification 中间件，检查是否实名制
func Certification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var whiteList = []string{"/docs", "/wx/user/login", "/ping", "/wx/user/get_user_info", "/wx/user/update_user_info", "/wx/info/project_types", "/wx/info/worker_types", "/wx/work/search"}

		var requestURL = c.Request.RequestURI
		for _, v := range whiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Next()
				return
			}
		}

		claims, exist := c.Get("claims")
		// 获取数据失败
		if common.FuncHandler(c, exist, true, common.SystemError) {
			return
		}
		userID := claims.(*model.CustomClaims).UserID
		user, _ := storage.UserExist(userID)
		realName := user.RealName
		phone := user.Phone

		realName = strings.Replace(realName, " ", "", -1)
		phone = strings.Replace(phone, " ", "", -1)

		if common.FuncHandler(c, realName == "" || phone == "", false, common.NoCertification) {
			c.Abort()
			return
		}
		c.Next()
	}
}
