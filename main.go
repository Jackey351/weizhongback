package main

import (
	"./common"
	"./controller/misc"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(common.ErrorHandling())

	r.GET("/ping", misc.Ping)
	r.Run() // listen and serve on 0.0.0.0:8080
}
