package wx

import (
	"fmt"
	"net/http"
	"time"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

const (
	// HourRecord 代表工时记录的常量
	HourRecord = 0
)

// AddHourRecord 添加工时记录
// @Summary 添加工时记录
// @Description 添加工时记录
// @Tags 工时
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

	fmt.Println(hourRecord.ID)
	err := tx.Create(&hourRecord).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}
	fmt.Println(hourRecord.ID)

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
