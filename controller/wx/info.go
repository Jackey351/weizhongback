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
// @Tags 各种类型信息
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/info/worker_types [get]
func GetWokerType(c *gin.Context) {
	var types []model.WorkType

	db := common.GetMySQL()

	err := db.Find(&types).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		return
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: types,
	})
}

// GetProjectType 获取所有工程类别
// @Summary 获取所有工程类别
// @Description 获取所有工程类别
// @Tags 各种类型信息
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/info/project_types [get]
func GetProjectType(c *gin.Context) {
	var types []model.ProjectType

	db := common.GetMySQL()

	err := db.Find(&types).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		return
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: types,
	})
}
