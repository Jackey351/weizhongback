package wx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// UserExist 根据用户id判断用户是否存在，不存在直接返回UserNoExist
func UserExist(c *gin.Context, userID int64) interface{} {
	db := common.GetMySQL()

	var existUser model.WxUser
	err := db.First(&existUser, userID).Error
	if common.FuncHandler(c, err, nil, common.UserNoExist) {
		return false
	}
	return existUser
}

const (
	worker = 3
)

// UpdateInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags 用户相关
// @Param token header string true "token"
// @Param user body model.WxUserInfo true "用户个人信息"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/user/update_user_info [post]
func UpdateInfo(c *gin.Context) {
	var userInfo model.WxUserInfo
	if common.FuncHandler(c, c.BindJSON(&userInfo), nil, common.ParameterError) {
		return
	}
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}

	userID := claims.(*model.CustomClaims).UserID
	var user model.WxUser
	var ok bool
	if user, ok = UserExist(c, userID).(model.WxUser); !ok {
		return
	}

	db := common.GetMySQL()
	tx := db.Begin()

	updateData := map[string]interface{}{"nick_name": userInfo.NickName, "real_name": userInfo.RealName, "sex": userInfo.Sex, "hometown": userInfo.Hometown, "phone": userInfo.Phone, "update_time": time.Now().Unix()}

	err := db.Model(&user).Updates(updateData).Error
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		tx.Rollback()
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, controller.Message{
		Msg: "更新成功",
	})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户相关
// @Param token header string true "token"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/user/get_user_info [get]
func GetUserInfo(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}

	userID := claims.(*model.CustomClaims).UserID
	var user model.WxUser
	var ok bool
	if user, ok = UserExist(c, userID).(model.WxUser); !ok {
		return
	}

	userInfo := user.WxUserInfo

	c.JSON(http.StatusOK, controller.Message{
		Data: userInfo,
	})
}

// Login 小程序用户登录
// @Summary 小程序用户登录
// @Description 小程序用户登录
// @Tags 用户相关
// @Param code query string true "登录码"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/user/login [get]
func Login(c *gin.Context) {
	code := c.Query("code")

	appID := viper.GetString("wechat.app_id")
	appSecret := viper.GetString("wechat.app_secret")

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, code)

	reponse, err := http.Get(url)
	if common.FuncHandler(c, err, nil, common.SystemError) {
		return
	}

	var data map[string]interface{}
	body, _ := ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &data)

	openID, exist := data["openid"]
	sessionKey := data["session_key"]

	if common.FuncHandler(c, exist, true, common.InvalidLogin) {
		return
	}

	var ret map[string]interface{}
	ret = make(map[string]interface{})
	var token string
	var userID int64

	// 利用openID搜索是否已存在，存在则更新，不存在插入新记录
	db := common.GetMySQL()
	tx := db.Begin()

	var existUser model.WxUser
	err = db.Where("open_id = ?", openID).First(&existUser).Error
	// 已有用户
	if err == nil {
		userID = existUser.ID

		var updateData = map[string]interface{}{"session_key": sessionKey.(string), "update_time": time.Now().Unix()}
		err := db.Model(&existUser).Updates(updateData).Error
		if common.FuncHandler(c, err, nil, common.DatabaseError) {
			tx.Rollback()
			return
		}

		token, err = common.CreateToken(userID)
		if common.FuncHandler(c, err, nil, common.SystemError) {
			return
		}

		tx.Commit()
	} else {
		var newUser model.WxUser
		newUser.OpenID = openID.(string)
		newUser.SessionKey = sessionKey.(string)
		newUser.Role = worker
		newUser.UpdateTime = time.Now().Unix()

		err := tx.Create(&newUser).Error
		if common.FuncHandler(c, err, nil, common.DatabaseError) {
			tx.Rollback()
			return
		}

		userID = newUser.ID
		token, err = common.CreateToken(userID)
		if common.FuncHandler(c, err, nil, common.SystemError) {
			tx.Rollback()
			return
		}

		tx.Commit()
	}

	ret["token"] = token
	ret["user_id"] = userID
	c.JSON(http.StatusOK, controller.Message{
		Data: ret,
	})
}
