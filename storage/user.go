package storage

import (
	"hackthoon/common"
	"hackthoon/model"
)

// UserExist 根据用户id判断用户是否存在，不存在直接返回UserNoExist
func UserExist(userID int64) (model.WxUser, error) {
	existUser, err := GetUserByID(userID)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

// UserPrefix 在区块链上的用户前缀
const UserPrefix = "resource:org.record.Worker#"

func getUserByIDFromDatabase(userID int64) (model.WxUser, error) {
	var existUser model.WxUser
	db := common.GetMySQL()
	err := db.First(&existUser, userID).Error
	return existUser, err
}

// GetUserByID 根据user_id获取用户
func GetUserByID(userID int64) (model.WxUser, error) {
	return getUserByIDFromDatabase(userID)
}

type user struct {
	Class    string `json:"$class"`
	UserID   string `json:"userId"`
	RealName string `json:"real_name"`
	Sex      string `json:"sex"`
	Hometown string `json:"hometown"`
	Phone    string `json:"phone"`
	NickName string `json:"nick_name"`
}

// SaveNewUser 保存新用户
func SaveNewUser(openID string, sessionKey string, role int, updateTime int64) (model.WxUser, error) {
	var newUser model.WxUser
	newUser.OpenID = openID
	newUser.SessionKey = sessionKey
	newUser.Role = role
	newUser.UpdateTime = updateTime

	db := common.GetMySQL()
	tx := db.Begin()

	err := tx.Create(&newUser).Error
	if err != nil {
		tx.Rollback()
		return newUser, err
	}

	tx.Commit()
	return newUser, nil
}

type userHua struct {
	Class    string `json:"$class"`
	UserID   int64  `json:"userId"`
	RealName string `json:"real_name"`
	Sex      string `json:"sex"`
	Hometown string `json:"hometown"`
	Phone    string `json:"phone"`
	NickName string `json:"nick_name"`
}

// UpdateUserInfo 更新用户
func UpdateUserInfo(user model.WxUser, nickName string, realName string, sex string, hometown string, phone string, updateTime int64) error {
	db := common.GetMySQL()
	tx := db.Begin()

	updateData := map[string]interface{}{"nick_name": nickName, "real_name": realName, "sex": sex, "hometown": hometown, "phone": phone, "update_time": updateTime}

	err := tx.Model(&user).Updates(updateData).Error
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
