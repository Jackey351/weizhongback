package storage

import (
	"yanfei_backend/common"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

// GroupExistByID 根据班组id判断班组是否存在，不存在直接返回GroupNoExist
func GroupExistByID(c *gin.Context, groupID int64) interface{} {
	db := common.GetMySQL()

	var existGroup model.Group
	err := db.First(&existGroup, groupID).Error
	if common.FuncHandler(c, err, nil, common.GroupNoExist) {
		return false
	}

	return existGroup
}
