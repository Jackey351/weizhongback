package main

import (
	"io"

	"yanfei_backend/common"
	"yanfei_backend/controller/misc"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func migrate(db *gorm.DB) {
	// 后面可以使用AutoMigrate，保持数据库的统一
	// AutoMigration只会根据struct tag建立新表、没有的列以及索引
	// 不会改变已经存在的列的类型或者删除没有用到的列
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1")
}

// init 在 main 之前执行
func init() {
	// init config
	common.DefaultConfig()
	common.SetConfig()
	common.WatchConfig()

	// init logger
	common.InitLogger()

	// init Database
	db := common.InitMySQL()
	migrate(db)
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
	r.Run()
}
