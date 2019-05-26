package wx

import (
	"net/http"
	"regexp"
	"strconv"
	"time"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"
	"yanfei_backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
// @Param token header string true "token"
// @Param 工时记录数据 body model.HourRecordRequest true "工时记录数据"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/add_hour_record [post]
func AddHourRecord(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var hourRecordRequest model.HourRecordRequest
	if common.FuncHandler(c, c.BindJSON(&hourRecordRequest), nil, common.ParameterError) {
		return
	}

	if _, ok := storage.UserExist(c, hourRecordRequest.WorkerID).(model.WxUser); !ok {
		return
	}
	if _, ok := storage.UserExist(c, userID).(model.WxUser); !ok {
		return
	}
	if _, ok := storage.GroupExistByID(c, hourRecordRequest.GroupID).(model.Group); !ok {
		return
	}

	var hourRecord model.HourRecord
	hourRecord.WorkHours = hourRecordRequest.WorkHours
	hourRecord.ExtraWorkHours = hourRecordRequest.ExtraWorkHours

	db := common.GetMySQL()

	var existRecord model.Record
	err := db.Where("worker_id = ? AND record_date = ?", hourRecordRequest.GroupID, hourRecordRequest.WorkerID, hourRecordRequest.RecordDate).First(&existRecord).Error
	if common.FuncHandler(c, err != nil, true, common.RecordHasExist) {
		return
	}

	tx := db.Begin()

	err = tx.Create(&hourRecord).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var record model.Record
	record.CommonRecord = hourRecordRequest.CommonRecord
	record.AdderID = userID
	record.RecordType = HourRecord
	record.RecordID = hourRecord.ID
	record.AddTime = time.Now().Unix()

	err = tx.Create(&record).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
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
// @Param token header string true "token"
// @Param 分项记录数据 body model.ItemRecordRequest true "分项记录数据"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/add_item_record [post]
func AddItemRecord(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var itemRecordRequest model.ItemRecordRequest
	if common.FuncHandler(c, c.BindJSON(&itemRecordRequest), nil, common.ParameterError) {
		return
	}

	if _, ok := storage.UserExist(c, itemRecordRequest.WorkerID).(model.WxUser); !ok {
		return
	}
	if _, ok := storage.UserExist(c, userID).(model.WxUser); !ok {
		return
	}
	if _, ok := storage.GroupExistByID(c, itemRecordRequest.GroupID).(model.Group); !ok {
		return
	}

	var itemRecord model.ItemRecord
	itemRecord.Subitem = itemRecordRequest.Subitem
	itemRecord.Quantity = itemRecordRequest.Quantity
	itemRecord.Unit = itemRecordRequest.Unit

	db := common.GetMySQL()

	var existRecord model.Record
	err := db.Where("worker_id = ? AND record_date = ?", itemRecordRequest.GroupID, itemRecordRequest.WorkerID, itemRecordRequest.RecordDate).First(&existRecord).Error
	if common.FuncHandler(c, err != nil, true, common.RecordHasExist) {
		return
	}

	tx := db.Begin()

	err = tx.Create(&itemRecord).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var record model.Record
	record.CommonRecord = itemRecordRequest.CommonRecord
	record.AdderID = userID
	record.RecordType = ItemRecord
	record.RecordID = itemRecord.ID
	record.AddTime = time.Now().Unix()

	err = tx.Create(&record).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
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
// @Param token header string true "token"
// @Param worker_id query int true "工人id"
// @Param date query string true "日期"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/check_recorded [get]
func CheckRecorded(c *gin.Context) {
	_, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}

	var workerID int64
	var err error

	workerID, err = strconv.ParseInt(c.Query("worker_id"), 10, 64)
	if common.FuncHandler(c, err, nil, common.ParameterError) {
		return
	}
	date := c.Query("date")

	if _, ok := storage.UserExist(c, workerID).(model.WxUser); !ok {
		return
	}

	var record model.Record
	db := common.GetMySQL()

	err = db.Where("worker_id = ? AND record_date = ?", workerID, date).First(&record).Error

	if err == nil {
		switch record.RecordType {
		case HourRecord:
			var hourRecordRequest model.HourRecordRequest
			hourRecordRequest.CommonRecord = record.CommonRecord

			var hourRecord model.HourRecord
			err = db.First(&hourRecord, record.RecordID).Error
			if common.FuncHandler(c, err, nil, common.DatabaseError) {
				return
			}

			hourRecordRequest.WorkHours = hourRecord.WorkHours
			hourRecordRequest.ExtraWorkHours = hourRecord.ExtraWorkHours

			var adderUser model.WxUser
			var workerUser model.WxUser
			var group model.Group
			var ok bool
			if adderUser, ok = storage.UserExist(c, record.AdderID).(model.WxUser); !ok {
				return
			}
			if workerUser, ok = storage.UserExist(c, workerID).(model.WxUser); !ok {
				return
			}
			if group, ok = storage.GroupExistByID(c, record.GroupID).(model.Group); !ok {
				return
			}

			var retHourInfo model.RetHourInfo
			retHourInfo.RecordID = record.ID
			retHourInfo.AdderInfo = adderUser.WxUserInfo
			retHourInfo.WorkerInfo = workerUser.WxUserInfo
			retHourInfo.GroupInfo = group.GroupRequest
			retHourInfo.RecordDate = record.RecordDate
			retHourInfo.Remark = record.Remark
			retHourInfo.WorkHours = hourRecord.WorkHours
			retHourInfo.ExtraWorkHours = hourRecord.ExtraWorkHours
			retHourInfo.AddTime = record.AddTime
			retHourInfo.IsConfirm = record.IsConfirm
			c.JSON(http.StatusOK, controller.Message{
				Data: retHourInfo,
			})
			break
		case ItemRecord:
			var itemRecordRequest model.ItemRecordRequest
			itemRecordRequest.CommonRecord = record.CommonRecord

			var itemRecord model.ItemRecord
			err = db.First(&itemRecord, record.RecordID).Error
			if common.FuncHandler(c, err, nil, common.DatabaseError) {
				return
			}

			itemRecordRequest.Subitem = itemRecord.Subitem
			itemRecordRequest.Quantity = itemRecord.Quantity

			var adderUser model.WxUser
			var workerUser model.WxUser
			var group model.Group
			var ok bool
			if adderUser, ok = storage.UserExist(c, record.AdderID).(model.WxUser); !ok {
				return
			}
			if workerUser, ok = storage.UserExist(c, workerID).(model.WxUser); !ok {
				return
			}
			if group, ok = storage.GroupExistByID(c, record.GroupID).(model.Group); !ok {
				return
			}

			var retItemInfo model.RetItemInfo
			retItemInfo.RecordID = record.ID
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
			c.JSON(http.StatusOK, controller.Message{
				Data: retItemInfo,
			})
			break
		}

	} else {
		c.JSON(http.StatusOK, controller.Message{
			Msg: "无记录",
		})
	}
}

// GetMonthRecords 查看某月的工作记录
// @Summary 查看某月的工作记录
// @Description 查看某月的工作记录
// @Tags 工作记录相关
// @Param token header string true "token"
// @Param month query string true "月份，形如2019-04"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/get_month_records [get]
func GetMonthRecords(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	month := c.Query("month")
	match, _ := regexp.MatchString("\\d{4}-\\d{2}", month)
	if common.FuncHandler(c, len(month) == 7 && match, true, common.ParameterError) {
		return
	}

	returnRecords, err := storage.GetRecordByMonth(c, userID, month)
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		return
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: returnRecords,
	})
}

// ConfirmRecord 确认工作记录
// @Summary 确认工作记录
// @Description 确认工作记录
// @Tags 工作记录相关
// @Param token header string true "token"
// @Param record_id query int true "工作记录id"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/record/confirm_record [get]
func ConfirmRecord(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	recordID := c.Query("record_id")

	db := common.GetMySQL()

	// 检查record_id 有效性，是否存在是否已确认
	var record model.Record
	err := db.First(&record, recordID).Error
	if common.FuncHandler(c, err, nil, common.RecordNoExist) {
		return
	}
	if common.FuncHandler(c, record.IsConfirm == 0, true, common.RecordHasConfirm) {
		return
	}

	// 检查userID权限，必须是相关方且不能确认自己提起的
	var group model.Group
	err = db.First(&group, record.GroupID).Error
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		return
	}

	groupOwnerID := group.OwnerID
	adderID := record.AdderID
	workerID := record.WorkerID
	if adderID == workerID {
		// 班组长确认
		if common.FuncHandler(c, userID == groupOwnerID, true, common.NoConfirmPermission) {
			return
		}
	} else {
		// 工人确认
		if common.FuncHandler(c, userID == workerID, true, common.NoConfirmPermission) {
			return
		}
	}

	tx := db.Begin()

	updateData := map[string]interface{}{"is_confirm": 1, "confirm_time": time.Now().Unix()}
	err = db.Model(&record).Updates(updateData).Error
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		tx.Rollback()
		return
	}

	if viper.GetString("basic.method") != "database" {
		switch record.RecordType {
		case HourRecord:
			err := storage.AddNewHourRecord(record)
			if common.FuncHandler(c, err, nil, common.BlockchainError) {
				tx.Rollback()
				return
			}
			break
		case ItemRecord:
			err := storage.AddNewItemRecord(record)
			if common.FuncHandler(c, err, nil, common.BlockchainError) {
				tx.Rollback()
				return
			}
			break
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, controller.Message{
		Msg: "确认成功",
	})
}
