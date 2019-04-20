package common

import (
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// MaintenanceHandling 维护模式中间件
func MaintenanceHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		if viper.GetBool("basic.maintenance") {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"err_code": 10008,
				"message":  Errors[10008],
			})
			log.Println(c.ClientIP(), "Maintenance mode is on")
			raven.CaptureMessage("Maintenance mode is on", map[string]string{"type": "maintenance"})
			c.Abort()
		}

		c.Next()
	}
}
