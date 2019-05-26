package storage

import (
	"yanfei_backend/common"
	"yanfei_backend/model"
)

// GroupExistByID 根据班组id判断班组是否存在，不存在直接返回GroupNoExist
func GroupExistByID(groupID int64) (model.Group, error) {
	db := common.GetMySQL()

	var existGroup model.Group
	err := db.First(&existGroup, groupID).Error
	if err != nil {
		return existGroup, err
	}

	return existGroup, nil
}
