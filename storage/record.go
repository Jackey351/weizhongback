package storage

import (
	"fmt"
	"hackthoon/common"
	"hackthoon/model"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
)

func execComdand(name string, arg ...string) error {
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	// cmd := exec.Command("/home/scy/.conda/envs/fisco/bin/python3.6", "console.py", "call", "TableTest", "0xafabb3d842b93a244096a803570931723c46d62b", "select", "fruit")
	cmd := exec.Command(name, arg...)
	fmt.Println(cmd)
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		return err
	}
	// 读取输出结果
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	log.Println(string(out))
	return nil
}

// AddNewHourRecord 添加新的工时记录
func AddNewHourRecord(record model.Record) error {
	db := common.GetMySQL()
	var hourRecord model.HourRecord
	db.First(&hourRecord, record.RecordID)

	err := execComdand("/home/scy/.conda/envs/fisco/bin/python3.6", "console.py", "sendtx", "Record", "0x2114e331011bf61d2c4fc8e9899495e8c2f58716", "insert_hour_record", strconv.FormatInt(record.RecordID, 10), strconv.FormatInt(record.AdderID, 10), strconv.FormatInt(record.GroupID, 10), record.RecordDate, strconv.FormatInt(record.AddTime, 10), strconv.FormatInt(int64(hourRecord.WorkHours), 10), strconv.FormatInt(int64(hourRecord.ExtraWorkHours), 10), strconv.FormatInt(record.ConfirmTime, 10), strconv.FormatInt(record.WorkerID, 10))

	return err
}

// AddNewItemRecord 添加新的工项记录
func AddNewItemRecord(record model.Record) error {
	db := common.GetMySQL()
	var itemRecord model.ItemRecord
	db.First(&itemRecord, record.RecordID)

	err := execComdand("/home/scy/.conda/envs/fisco/bin/python3.6", "console.py", "sendtx", "Record", "0x2114e331011bf61d2c4fc8e9899495e8c2f58716", "insert_item_record", strconv.FormatInt(record.RecordID, 10), strconv.FormatInt(record.AdderID, 10), strconv.FormatInt(record.GroupID, 10), record.RecordDate, strconv.FormatInt(record.AddTime, 10), itemRecord.Subitem, strconv.FormatInt(int64(itemRecord.Quantity), 10), strconv.FormatInt(record.ConfirmTime, 10), strconv.FormatInt(record.WorkerID, 10))

	return err
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

// GetRecordByMonth 获取一个月的工作记录
func GetRecordByMonth(userID int64, month string) ([]interface{}, error) {
	return getRecordByMonthFromDatabase(userID, month)
}
