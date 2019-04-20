package controller

import (
	"net/http"
	"yanfei_backend/common"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	password := c.DefaultQuery("password", "")

	if common.FuncHandler(c, (name != "" && password != ""), true, 20000) {
		return
	}

	db := common.GetMySQL()

	var user model.User
	db.Where(&model.User{Name: name, Password: password}).First(&user)

	if common.FuncHandler(c, user.Name != "", true, 20000) {
		return
	}

	c.JSON(http.StatusOK, user)
}
