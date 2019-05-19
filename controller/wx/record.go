package wx

import (
	"net/http"
	"strconv"
	"time"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

// 代表不同的记录类型
const (
	HourRecord = 0
	ItemRecord = 1
)

// AddHourRecord 添加工时记录
// @Summary 添加工时记录
// @Description 添加工时记录
// @Tags 工作记录相关
// @Param 工时记录数据 body model.HourRecordRequest true "工时记录数据"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/add_hour_record [post]
func AddHourRecord(c *gin.Context) {
	var hourRecordRequest model.HourRecordRequest
	if common.FuncHandler(c, c.BindJSON(&hourRecordRequest), nil, 20001) {
		return
	}

	var hourRecord model.HourRecord
	hourRecord.WorkHours = hourRecordRequest.WorkHours
	hourRecord.ExtraWorkHours = hourRecordRequest.ExtraWorkHours

	db := common.GetMySQL()
	tx := db.Begin()

	err := tx.Create(&hourRecord).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var record model.Record
	record.CommonRecord = hourRecordRequest.CommonRecord
	record.RecordType = HourRecord
	record.RecordID = hourRecord.ID
	record.AddTime = time.Now().Unix()

	err = tx.Create(&record).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Msg: "添加成功",
	})
}

// AddItemRecord 添加分项记录
// @Summary 添加分项记录
// @Description 添加分项记录
// @Tags 工作记录相关
// @Param 分项记录数据 body model.ItemRecordRequest true "分项记录数据"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/add_item_record [post]
func AddItemRecord(c *gin.Context) {
	var itemRecordRequest model.ItemRecordRequest
	if common.FuncHandler(c, c.BindJSON(&itemRecordRequest), nil, 20001) {
		return
	}

	var itemRecord model.ItemRecord
	itemRecord.Subitem = itemRecordRequest.Subitem
	itemRecord.Quantity = itemRecordRequest.Quantity

	db := common.GetMySQL()
	tx := db.Begin()

	err := tx.Create(&itemRecord).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var record model.Record
	record.CommonRecord = itemRecordRequest.CommonRecord
	record.RecordType = ItemRecord
	record.RecordID = itemRecord.ID
	record.AddTime = time.Now().Unix()

	err = tx.Create(&record).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Msg: "添加成功",
	})
}

// CheckRecorded 检查某日是否记录
// @Summary 检查某日是否记录
// @Description 检查某日是否记录
// @Tags 工作记录相关
// @Param group_id query int true "班组id"
// @Param worker_id query int true "工人id"
// @Param date query string true "日期"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/check_recorded [get]
func CheckRecorded(c *gin.Context) {
	var groupID int64
	var workerID int64
	var err error

	groupID, err = strconv.ParseInt(c.Query("group_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}
	workerID, err = strconv.ParseInt(c.Query("worker_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}
	date := c.Query("date")

	var record model.Record
	db := common.GetMySQL()

	err = db.Where("group_id = ? AND worker_id = ? AND record_date = ?", groupID, workerID, date).First(&record).Error

	if err == nil {
		switch record.RecordType {
		case HourRecord:
			var hourRecordRequest model.HourRecordRequest
			hourRecordRequest.CommonRecord = record.CommonRecord

			var hourRecord model.HourRecord
			err = db.First(&hourRecord, record.RecordID).Error
			if common.FuncHandler(c, err, nil, 20002) {
				return
			}

			hourRecordRequest.WorkHours = hourRecord.WorkHours
			hourRecordRequest.ExtraWorkHours = hourRecord.ExtraWorkHours
			c.JSON(http.StatusOK, controller.Message{
				Data: hourRecordRequest,
			})
			break
		case ItemRecord:
			var itemRecordRequest model.ItemRecordRequest
			itemRecordRequest.CommonRecord = record.CommonRecord

			var itemRecord model.ItemRecord
			err = db.First(&itemRecord, record.RecordID).Error
			if common.FuncHandler(c, err, nil, 20002) {
				return
			}

			itemRecordRequest.Subitem = itemRecord.Subitem
			itemRecordRequest.Quantity = itemRecord.Quantity
			c.JSON(http.StatusOK, controller.Message{
				Data: itemRecordRequest,
			})
			break

		}
	} else {
		c.JSON(http.StatusOK, controller.Message{
			Msg: "无记录",
		})
	}
}
