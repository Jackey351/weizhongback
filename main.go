package main

import (
	"io"

	"hackthoon/common"
	"hackthoon/controller"
	"hackthoon/controller/wx"
	"hackthoon/middleware"

	_ "hackthoon/docs"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
	// 禁止在表名后面加s
	db.SingularTable(true)
	migrate(db)
}

// @title YANFEI API
// @version 0.0.1
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
	// middleware
	r.Use(middleware.ErrorHandling())
	r.Use(middleware.MaintenanceHandling())
	r.Use(middleware.JWTAuth())
	r.Use(middleware.Certification())

	// swagger router
	if viper.GetBool("basic.debug") {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 路由
	r.GET("/ping", controller.Ping)
	// 用户相关
	r.POST("/wx/user/update_user_info", wx.UpdateInfo)
	r.GET("/wx/user/login", wx.Login)
	r.GET("/wx/user/get_user_info", wx.GetUserInfo)
	// 各种类型信息
	r.GET("/wx/info/worker_types", wx.GetWokerType)
	r.GET("/wx/info/project_types", wx.GetProjectType)
	// 工作相关
	r.POST("/wx/work/publish_dian", wx.PublishDianWork)
	r.POST("/wx/work/publish_bao", wx.PublishBaoWork)
	r.POST("/wx/work/publish_tuji", wx.PublishTujiWork)
	r.GET("/wx/work/search", wx.SearchWork)
	// 班组相关
	r.POST("/wx/group/new_group", wx.NewGroup)
	r.GET("/wx/group/join_group", wx.JoinGroup)
	r.GET("/wx/group/in_group", wx.InGroup)
	r.GET("/wx/group/group_member", wx.GroupMember)
	r.GET("/wx/group/delete_member", wx.DeleteMember)
	// 工作记录相关
	r.POST("/wx/record/add_hour_record", wx.AddHourRecord)
	r.POST("/wx/record/add_item_record", wx.AddItemRecord)
	r.GET("/wx/record/check_recorded", wx.CheckRecorded)
	r.GET("/wx/record/get_month_records", wx.GetMonthRecords)
	r.GET("/wx/record/confirm_record", wx.ConfirmRecord)

	r.Run("0.0.0.0:" + viper.GetString("basic.port"))
}
