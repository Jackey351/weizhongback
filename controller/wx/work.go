package wx

import (
	"net/http"
	"strings"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

// PublishWork 发布工作
// @Summary 发布工作
// @Description 发布工作
// @Tags wx
// @Param type query string true "工种 0(点工),1(包工)"
// @Param 点工示例数据 body model.DianWorkWrapper false "点工招聘"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/work/publish [post]
func PublishWork(c *gin.Context) {
	priceType := c.Query("type")

	if common.FuncHandler(c, priceType == "0" || priceType == "1", true, 20001) {
		return
	}

	if priceType == "0" {
		var workWrapper model.DianWorkWrapper
		if common.FuncHandler(c, c.BindJSON(&workWrapper), nil, 20001) {
			return
		}

		db := common.GetMySQL()

		// 检查工种和工程类别是否正确
		projectType := workWrapper.BasicWork.ProjectType
		workerType := workWrapper.BasicWork.WorkerType

		var res model.ProjectType
		err := db.Where("name = ?", projectType).First(&res).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, 30000) {
			return
		}

		var res2 model.WorkerType
		err = db.Where("name = ?", workerType).First(&res2).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, 30001) {
			return
		}

		var dianWork model.DianWork
		dianWork.DianWorkOther = workWrapper.DianWorkOther

		tx := db.Begin()

		err = tx.Create(&dianWork).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		var locationInfo model.LocationInfo
		locationInfo.LocationInfoWrapper = workWrapper.WorkWrapper.LocationInfoWrapper
		err = tx.Create(&locationInfo).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		var work model.Work
		work.BasicWork = workWrapper.WorkWrapper.BasicWork
		work.LocationID = locationInfo.ID
		work.Treatment = strings.Join(workWrapper.WorkWrapper.Treatment, ", ")
		work.Fid = dianWork.ID

		err = tx.Create(&work).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, controller.Message{
			Data: "发布成功",
		})
	} else {
		var workWrapper model.BaoWorkWrapper
		if common.FuncHandler(c, c.BindJSON(&workWrapper), nil, 20001) {
			return
		}

		db := common.GetMySQL()

		// 检查工种和工程类别是否正确
		projectType := workWrapper.BasicWork.ProjectType
		workerType := workWrapper.BasicWork.WorkerType

		var res model.ProjectType
		err := db.Where("name = ?", projectType).First(&res).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, 30000) {
			return
		}

		var res2 model.WorkerType
		err = db.Where("name = ?", workerType).First(&res2).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, 30001) {
			return
		}

		var baoWork model.BaoWork
		baoWork.BaoWorkOther = workWrapper.BaoWorkOther

		tx := db.Begin()

		err = tx.Create(&baoWork).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		var locationInfo model.LocationInfo
		locationInfo.LocationInfoWrapper = workWrapper.WorkWrapper.LocationInfoWrapper
		err = tx.Create(&locationInfo).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		var work model.Work
		work.BasicWork = workWrapper.WorkWrapper.BasicWork
		work.LocationID = locationInfo.ID
		work.Treatment = strings.Join(workWrapper.WorkWrapper.Treatment, ", ")
		work.Fid = baoWork.ID

		err = tx.Create(&work).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, controller.Message{
			Data: "发布成功",
		})
	}

}
