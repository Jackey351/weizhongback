package main

import (
	"io"

	"./common"
	"./controller/misc"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

// init 在 main 之前执行
func init() {
	// init config
	common.DefaultConfig()
	common.SetConfig()
	common.WatchConfig()

	// init logger
	common.InitLogger()
}

func main() {
	// Before init router
	if viper.GetBool("basic.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		// Redirect log to file
		gin.DisableConsoleColor()
		logFile := common.GetLogFile()
		defer logFile.Close()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	r := gin.Default()
	// Error handling
	r.Use(common.ErrorHandling())
	r.Use(common.MaintenanceHandling())

	r.GET("/ping", misc.Ping)
	r.Run() // listen and serve on 0.0.0.0:8080
}
