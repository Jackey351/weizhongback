package wx

import (
	"net/http"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

// NewWxUser 小程序新用户
// @Summary 小程序端新添用户
// @Description 小程序端新添用户
// @Tags 用户相关
// @Param user body model.WxUserWrapper true "create a new user"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/user/new_user [post]
func NewWxUser(c *gin.Context) {
	var newUserReq model.WxUserWrapper

	// 获取数据失败
	if common.FuncHandler(c, c.BindJSON(&newUserReq), nil, common.ParameterError) {
		return
	}

	var newUser model.WxUser
	newUser.WxUserWrapper = newUserReq
	newUser.Role = 3

	db := common.GetMySQL()
	tx := db.Begin()

	err := tx.Create(&newUser).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, controller.Message{
		Data: newUser,
	})
}
