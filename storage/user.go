package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"yanfei_backend/common"
	"yanfei_backend/model"

	"github.com/spf13/viper"
)

// UserPrefix 在区块链上的用户前缀
const UserPrefix = "resource:org.record.Worker#"

func getUserByIDFromDatabase(userID int64) (model.WxUser, error) {
	var existUser model.WxUser
	db := common.GetMySQL()
	err := db.First(&existUser, userID).Error
	return existUser, err
}

func getUserByIDFromHyperledger(userID int64) (model.WxUser, error) {
	var existUser model.WxUser
	db := common.GetMySQL()
	err := db.First(&existUser, userID).Error

	if err != nil {
		return existUser, err
	}

	basicURL := viper.GetString("blockchain.hyperledger.url")
	getUserInfoAPI := fmt.Sprintf("%s/org.record.Worker/%d", basicURL, userID)
	reponse, err := http.Get(getUserInfoAPI)

	if err != nil {
		return existUser, err
	}
	if reponse.StatusCode != 200 {
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(reponse.Body)
		json.Unmarshal(body, &data)
		fmt.Println(data)
		return existUser, errors.New("系统出错")
	}

	var data map[string]interface{}
	body, _ := ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &data)

	existUser.Hometown = data["hometown"].(string)
	existUser.RealName = data["real_name"].(string)
	existUser.Sex = data["sex"].(string)
	existUser.Phone = data["phone"].(string)
	existUser.NickName = data["nick_name"].(string)

	return existUser, err
}

// GetUserByID 根据user_id获取用户
func GetUserByID(userID int64) (model.WxUser, error) {
	switch viper.GetString("basic.method") {
	default:
		return getUserByIDFromDatabase(userID)
	case "hyperledger":
		return getUserByIDFromHyperledger(userID)
	}
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

func saveNewUserToHyperledger(userID int64) error {
	basicURL := viper.GetString("blockchain.hyperledger.url")
	newUserAPI := fmt.Sprintf("%s/org.record.Worker", basicURL)

	var newUser user
	newUser.Class = "org.record.Worker"
	newUser.UserID = strconv.FormatInt(userID, 10)
	newUser.RealName = " "
	newUser.Sex = " "
	newUser.Hometown = " "
	newUser.Phone = " "
	newUser.NickName = " "

	b, err := json.Marshal(newUser)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer([]byte(b))
	reponse, err2 := http.Post(newUserAPI, "application/json;charset=utf-8", body)

	if err2 != nil {
		return err
	}
	if reponse.StatusCode != 200 {
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(reponse.Body)
		json.Unmarshal(body, &data)
		fmt.Println(data)
		return errors.New("系统出错")
	}

	return nil
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

	userID := newUser.ID

	switch viper.GetString("basic.method") {
	case "hyperledger":
		err = saveNewUserToHyperledger(userID)
		if err != nil {
			tx.Rollback()
			return newUser, err
		}
		break
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

func updateUserInfoToHyperledger(userID int64, nickName string, realName string, sex string, hometown string, phone string) error {
	basicURL := viper.GetString("blockchain.hyperledger.url")
	updateUserAPI := fmt.Sprintf("%s/org.record.Worker/%d", basicURL, userID)

	var newUser userHua
	newUser.Class = "org.record.Worker"
	newUser.UserID = userID
	newUser.RealName = realName
	newUser.Sex = sex
	newUser.Hometown = hometown
	newUser.Phone = phone
	newUser.NickName = nickName

	b, err := json.Marshal(newUser)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer([]byte(b))

	req, _ := http.NewRequest("PUT", updateUserAPI, body)
	req.Header.Add("Content-Type", "application/json")

	reponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if reponse.StatusCode != 200 {
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(reponse.Body)
		json.Unmarshal(body, &data)
		fmt.Println(data)
		return errors.New("系统出错")
	}

	return nil
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

	switch viper.GetString("basic.method") {
	case "hyperledger":
		err = updateUserInfoToHyperledger(user.ID, nickName, realName, sex, hometown, phone)
		if err != nil {
			tx.Rollback()
			return err
		}
		break
	}

	tx.Commit()
	return nil
}
