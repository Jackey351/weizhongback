package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary PING-PONG
// @Description 测试服务器是否在线
// @Tags 测试
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, Message{
		Msg: "pong",
	})
}
