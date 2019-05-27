package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"yanfei_backend/common"
	"yanfei_backend/model"

	"github.com/spf13/viper"
)

type hourRecordBC struct {
	Class       string  `json:"$class"`
	ID          string  `json:"worktimeId"`
	AdderID     int64   `json:"adder_id"`
	GroupID     int64   `json:"group_id"`
	Date        string  `json:"date"`
	AddTime     int64   `json:"add_time"`
	Hours       float64 `json:"work_hours"`
	ExtraHours  float64 `json:"extra_work_hours"`
	Remark      string  `json:"remark"`
	ConfirmTime int64   `json:"confirm_time"`
	Owner       string  `json:"owner"`
}

type itemRecordBC struct {
	Class       string  `json:"$class"`
	ID          string  `json:"itemworktimeId"`
	AdderID     int64   `json:"adder_id"`
	GroupID     int64   `json:"group_id"`
	Date        string  `json:"date"`
	AddTime     int64   `json:"add_time"`
	Subitem     string  `json:"subitem"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	Remark      string  `json:"remark"`
	ConfirmTime int64   `json:"confirm_time"`
	Owner       string  `json:"owner"`
}

// AddNewHourRecord 添加新的工时记录
func AddNewHourRecord(record model.Record) error {
	db := common.GetMySQL()
	var hourRecord model.HourRecord
	db.First(&hourRecord, record.RecordID)

	var hourRecordBC hourRecordBC
	hourRecordBC.Class = "org.record.Worktime"
	hourRecordBC.ID = strconv.FormatInt(record.RecordID, 10)
	hourRecordBC.AdderID = record.AdderID
	hourRecordBC.GroupID = record.GroupID
	hourRecordBC.Date = record.RecordDate
	hourRecordBC.AddTime = record.AddTime
	hourRecordBC.Hours = hourRecord.WorkHours
	hourRecordBC.ExtraHours = hourRecord.ExtraWorkHours
	hourRecordBC.Remark = record.Remark
	hourRecordBC.ConfirmTime = record.ConfirmTime
	hourRecordBC.Owner = fmt.Sprintf("%s%d", UserPrefix, record.WorkerID)

	b, err := json.Marshal(hourRecordBC)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer([]byte(b))

	basicURL := viper.GetString("blockchain.hyperledger.url")
	API := fmt.Sprintf("%s/org.record.Worktime", basicURL)
	reponse, err := http.Post(API, "application/json;charset=utf-8", body)

	if err != nil {
		return err
	}
	if reponse.StatusCode != 200 {
		return errors.New("系统出错")
	}

	return nil
}

// AddNewItemRecord 添加新的工项记录
func AddNewItemRecord(record model.Record) error {
	db := common.GetMySQL()
	var itemRecord model.ItemRecord
	db.First(&itemRecord, record.RecordID)

	var itemRecordBC itemRecordBC
	itemRecordBC.Class = "org.record.Itemworktime"
	itemRecordBC.ID = strconv.FormatInt(record.RecordID, 10)
	itemRecordBC.AdderID = record.AdderID
	itemRecordBC.GroupID = record.GroupID
	itemRecordBC.Date = record.RecordDate
	itemRecordBC.AddTime = record.AddTime
	itemRecordBC.Subitem = itemRecord.Subitem
	itemRecordBC.Quantity = itemRecord.Quantity
	itemRecordBC.Unit = itemRecord.Unit
	itemRecordBC.Remark = record.Remark
	itemRecordBC.ConfirmTime = record.ConfirmTime
	itemRecordBC.Owner = fmt.Sprintf("%s%d", UserPrefix, record.WorkerID)

	b, err := json.Marshal(itemRecordBC)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer([]byte(b))

	basicURL := viper.GetString("blockchain.hyperledger.url")
	API := fmt.Sprintf("%s/org.record.Itemworktime", basicURL)
	reponse, err := http.Post(API, "application/json;charset=utf-8", body)

	if err != nil {
		return err
	}
	if reponse.StatusCode != 200 {
		return errors.New("系统出错")
	}

	return nil
}

const (
	hourRecord = 0
	itemRecord = 1
)

func getRecordByMonthFromDatabase(userID int64, month string) ([]interface{}, error) {
	month = month + "%"
	db := common.GetMySQL()

	var records []model.Record
	err := db.Where("is_confirm = 1 AND worker_id = ? AND record_date LIKE ?", userID, month).Order("record_date asc").Find(&records).Error

	var returnRecords []interface{}
	if err == nil {
		for _, record := range records {
			switch record.RecordType {
			case hourRecord:
				var retHourInfo model.RetHourInfo

				var hourRecord model.HourRecord
				err = db.First(&hourRecord, record.RecordID).Error
				if err != nil {
					return returnRecords, err
				}

				var adderUser model.WxUser
				var workerUser model.WxUser
				var group model.Group

				adderUser, err = UserExist(record.AdderID)
				if err != nil {
					return returnRecords, err
				}
				workerUser, err = UserExist(record.WorkerID)
				if err != nil {
					return returnRecords, err
				}
				group, err = GroupExistByID(record.GroupID)
				if err != nil {
					return returnRecords, err
				}

				retHourInfo.RecordID = hourRecord.ID
				retHourInfo.AdderInfo = adderUser.WxUserInfo
				retHourInfo.WorkerInfo = workerUser.WxUserInfo
				retHourInfo.GroupInfo = group.GroupRequest
				retHourInfo.RecordDate = record.RecordDate
				retHourInfo.Remark = record.Remark
				retHourInfo.WorkHours = hourRecord.WorkHours
				retHourInfo.ExtraWorkHours = hourRecord.ExtraWorkHours
				retHourInfo.AddTime = record.AddTime
				retHourInfo.IsConfirm = record.IsConfirm
				retHourInfo.Type = 0

				returnRecords = append(returnRecords, retHourInfo)
				break
			case itemRecord:
				var itemRecordRequest model.ItemRecordRequest
				itemRecordRequest.CommonRecord = record.CommonRecord

				var itemRecord model.ItemRecord
				err = db.First(&itemRecord, record.RecordID).Error
				if err != nil {
					return returnRecords, err
				}

				var adderUser model.WxUser
				var workerUser model.WxUser
				var group model.Group

				adderUser, err = UserExist(record.AdderID)
				if err != nil {
					return returnRecords, err
				}
				workerUser, err = UserExist(record.WorkerID)
				if err != nil {
					return returnRecords, err
				}
				group, err = GroupExistByID(record.GroupID)
				if err != nil {
					return returnRecords, err
				}

				var retItemInfo model.RetItemInfo
				retItemInfo.RecordID = itemRecord.ID
				retItemInfo.AdderInfo = adderUser.WxUserInfo
				retItemInfo.WorkerInfo = workerUser.WxUserInfo
				retItemInfo.GroupInfo = group.GroupRequest
				retItemInfo.RecordDate = record.RecordDate
				retItemInfo.Remark = record.Remark
				retItemInfo.Subitem = itemRecord.Subitem
				retItemInfo.Quantity = itemRecord.Quantity
				retItemInfo.Unit = itemRecord.Unit
				retItemInfo.AddTime = record.AddTime
				retItemInfo.IsConfirm = record.IsConfirm
				retItemInfo.Type = 1

				returnRecords = append(returnRecords, retItemInfo)
				break
			}
		}
	}
	return returnRecords, nil

}

func getRecordByMonthFromHyperledger(userID int64, month string) ([]interface{}, error) {
	var returnRecords []interface{}

	var hourRecords []model.RetHourInfo
	var itemRecords []model.RetItemInfo

	regex := regexp.MustCompile(`^([\d]{4})-([\d]{2})$`)
	params := regex.FindStringSubmatch(month)

	minYear, _ := strconv.Atoi(params[1])
	minMonth, _ := strconv.Atoi(params[2])

	maxYear := minYear
	maxMonth := minMonth + 1
	if maxMonth > 12 {
		maxMonth = 1
		maxYear = maxYear + 1
	}

	basicURL := viper.GetString("blockchain.hyperledger.url")
	API := fmt.Sprintf("%s/queries/selectItemWorktimeByMonthandOwner", basicURL)

	p := fmt.Sprintf("?owner=%s%d&moins=%d-%02d&plus=%d-%02d", UserPrefix, userID, minYear, minMonth, maxYear, maxMonth)
	p = strings.Replace(p, ":", "%3A", -1)
	p = strings.Replace(p, "#", "%23", -1)

	API = API + p
	reponse, err := http.Get(API)

	if err != nil {
		return returnRecords, err
	}
	if reponse.StatusCode != 200 {
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(reponse.Body)
		json.Unmarshal(body, &data)
		fmt.Println(data)
		return returnRecords, errors.New("系统出错")
	}

	var datas []map[string]interface{}
	body, _ := ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &datas)

	for _, data := range datas {
		groupID := int64(data["group_id"].(float64))
		adderID := int64(data["adder_id"].(float64))
		itemworktimeID := data["itemworktimeId"].(string)
		recordID, _ := strconv.ParseInt(itemworktimeID, 10, 64)

		var adderUser model.WxUser
		var workerUser model.WxUser
		var group model.Group
		adderUser, err = UserExist(adderID)
		if err != nil {
			return returnRecords, err
		}
		workerUser, err = UserExist(userID)
		if err != nil {
			return returnRecords, err
		}
		group, err = GroupExistByID(groupID)
		if err != nil {
			return returnRecords, err
		}

		var retItemInfo model.RetItemInfo
		retItemInfo.RecordID = recordID
		retItemInfo.AdderInfo = adderUser.WxUserInfo
		retItemInfo.WorkerInfo = workerUser.WxUserInfo
		retItemInfo.GroupInfo = group.GroupRequest
		retItemInfo.RecordDate = data["date"].(string)
		retItemInfo.Remark = data["remark"].(string)
		retItemInfo.Subitem = data["subitem"].(string)
		retItemInfo.Quantity = data["quantity"].(float64)
		retItemInfo.Unit = data["unit"].(string)
		retItemInfo.IsConfirm = 1
		retItemInfo.AddTime = int64(data["add_time"].(float64))
		retItemInfo.Type = 1

		itemRecords = append(itemRecords, retItemInfo)
	}

	API = fmt.Sprintf("%s/queries/selectWorktimeByMonthandOwner", basicURL)

	p = fmt.Sprintf("?owner=%s%d&moins=%d-%02d&plus=%d-%02d", UserPrefix, userID, minYear, minMonth, maxYear, maxMonth)
	p = strings.Replace(p, ":", "%3A", -1)
	p = strings.Replace(p, "#", "%23", -1)

	API = API + p
	reponse, err = http.Get(API)

	if err != nil {
		return returnRecords, err
	}
	if reponse.StatusCode != 200 {
		var data map[string]interface{}
		body, _ := ioutil.ReadAll(reponse.Body)
		json.Unmarshal(body, &data)
		fmt.Println(data)
		return returnRecords, errors.New("系统出错")
	}

	body, _ = ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &datas)

	for _, data := range datas {
		groupID := int64(data["group_id"].(float64))
		adderID := int64(data["adder_id"].(float64))
		worktimeID := data["worktimeId"].(string)
		recordID, _ := strconv.ParseInt(worktimeID, 10, 64)

		var adderUser model.WxUser
		var workerUser model.WxUser
		var group model.Group
		adderUser, err = UserExist(adderID)
		if err != nil {
			return returnRecords, err
		}
		workerUser, err = UserExist(userID)
		if err != nil {
			return returnRecords, err
		}
		group, err = GroupExistByID(groupID)
		if err != nil {
			return returnRecords, err
		}

		var retHourInfo model.RetHourInfo
		retHourInfo.RecordID = recordID
		retHourInfo.AdderInfo = adderUser.WxUserInfo
		retHourInfo.WorkerInfo = workerUser.WxUserInfo
		retHourInfo.GroupInfo = group.GroupRequest
		retHourInfo.RecordDate = data["date"].(string)
		retHourInfo.Remark = data["remark"].(string)
		retHourInfo.WorkHours = data["work_hours"].(float64)
		retHourInfo.ExtraWorkHours = data["extra_work_hours"].(float64)
		retHourInfo.IsConfirm = 1
		retHourInfo.AddTime = int64(data["add_time"].(float64))
		retHourInfo.Type = 0

		hourRecords = append(hourRecords, retHourInfo)
	}

	i := 0
	j := 0

	for i != len(hourRecords) && j != len(itemRecords) {
		if hourRecords[i].RecordDate < itemRecords[j].RecordDate {
			returnRecords = append(returnRecords, hourRecords[i])
			i++
		} else {
			returnRecords = append(returnRecords, itemRecords[j])
			j++
		}
	}

	for i < len(hourRecords) {
		returnRecords = append(returnRecords, hourRecords[i])
		i++
	}
	for j < len(itemRecords) {
		returnRecords = append(returnRecords, itemRecords[j])
		j++
	}

	return returnRecords, nil
}

// GetRecordByMonth 获取一个月的工作记录
func GetRecordByMonth(userID int64, month string) ([]interface{}, error) {
	switch viper.GetString("basic.method") {
	default:
		return getRecordByMonthFromDatabase(userID, month)
	case "hyperledger":
		return getRecordByMonthFromHyperledger(userID, month)
	}
}
