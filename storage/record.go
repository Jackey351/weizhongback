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
				var hourRecordRequest model.HourRecordRequest
				hourRecordRequest.CommonRecord = record.CommonRecord

				var hourRecord model.HourRecord
				err = db.First(&hourRecord, record.RecordID).Error
				if err != nil {
					return returnRecords, err
				}

				hourRecordRequest.WorkHours = hourRecord.WorkHours
				hourRecordRequest.ExtraWorkHours = hourRecord.ExtraWorkHours

				returnRecords = append(returnRecords, hourRecordRequest)
				break
			case itemRecord:
				var itemRecordRequest model.ItemRecordRequest
				itemRecordRequest.CommonRecord = record.CommonRecord

				var itemRecord model.ItemRecord
				err = db.First(&itemRecord, record.RecordID).Error
				if err != nil {
					return returnRecords, err
				}

				itemRecordRequest.Subitem = itemRecord.Subitem
				itemRecordRequest.Quantity = itemRecord.Quantity

				returnRecords = append(returnRecords, itemRecordRequest)
				break
			}
		}
	}
	return returnRecords, nil

}

func getRecordByMonthFromHyperledger(userID int64, month string) ([]interface{}, error) {
	var returnRecords []interface{}

	var hourRecords []model.HourRecordRequest
	var itemRecords []model.ItemRecordRequest

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

	for _, param := range params {
		fmt.Println(param)
	}

	basicURL := viper.GetString("blockchain.hyperledger.url")
	API := fmt.Sprintf("%s/queries/selectItemWorktimeByMonthandOwner?owner=%s%d&moins=%d-%02d&plus=%d-%02d", basicURL, UserPrefix, userID, minYear, minMonth, maxYear, maxMonth)
	reponse, err := http.Get(API)

	if err != nil {
		return returnRecords, err
	}

	var datas []map[string]interface{}
	body, _ := ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &datas)

	for _, data := range datas {
		var itemRecord model.ItemRecordRequest
		itemRecord.GroupID = data["group_id"].(int64)
		itemRecord.Quantity = data["quantity"].(float64)
		itemRecord.RecordDate = data["date"].(string)
		itemRecord.Remark = data["remark"].(string)
		itemRecord.Subitem = data["subitem"].(string)
		itemRecord.Unit = data["unit"].(string)
		itemRecord.WorkerID, err = strconv.ParseInt(string([]byte(data["owner"].(string))[len(UserPrefix):]), 10, 64)

		if err != nil {
			return returnRecords, err
		}

		itemRecords = append(itemRecords, itemRecord)
	}

	API = fmt.Sprintf("%s/queries/selectWorktimeByMonthandOwner?owner=%s%d&moins=%d-%02d&plus=%d-%02d", basicURL, UserPrefix, userID, minYear, minMonth, maxYear, maxMonth)
	reponse, err = http.Get(API)

	if err != nil {
		return returnRecords, err
	}

	body, _ = ioutil.ReadAll(reponse.Body)
	json.Unmarshal(body, &datas)

	for _, data := range datas {
		var hourRecord model.HourRecordRequest
		hourRecord.GroupID = data["group_id"].(int64)
		hourRecord.WorkHours = data["work_hours"].(float64)
		hourRecord.RecordDate = data["date"].(string)
		hourRecord.Remark = data["remark"].(string)
		hourRecord.ExtraWorkHours = data["extra_work_hours"].(float64)
		hourRecord.WorkerID, err = strconv.ParseInt(string([]byte(data["owner"].(string))[len(UserPrefix):]), 10, 64)

		if err != nil {
			return returnRecords, err
		}

		hourRecords = append(hourRecords, hourRecord)
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
