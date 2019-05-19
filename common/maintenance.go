package common

import (
	"log"
	"net/http"
	"yanfei_backend/controller"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// MaintenanceHandling 维护模式中间件
func MaintenanceHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		if viper.GetBool("basic.maintenance") {
			c.JSON(http.StatusServiceUnavailable, controller.Message{
				Status: Maintenance,
				Msg:    Errors[Maintenance],
			})
			log.Println(c.ClientIP(), "Maintenance mode is on")
			raven.CaptureMessage("Maintenance mode is on", map[string]string{"type": "maintenance"})
			c.Abort()
		}

		c.Next()
	}
}
