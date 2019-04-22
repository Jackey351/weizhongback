package wx

import (
	"net/http"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

// GetWokerType 获取所有工种
// @Summary 获取所有工种
// @Description 获取所有工种
// @Tags wx
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/info/worker_types [get]
func GetWokerType(c *gin.Context) {
	var types []model.WorkerType

	db := common.GetMySQL()

	err := db.Find(&types).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		return
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: types,
	})
}

// GetProjectType 获取所有工程类别
// @Summary 获取所有工程类别
// @Description 获取所有工程类别
// @Tags wx
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/info/project_types [get]
func GetProjectType(c *gin.Context) {
	var types []model.ProjectType

	db := common.GetMySQL()

	err := db.Find(&types).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		return
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: types,
	})
}
