package wx

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

// PublishWork 发布工作
// @Summary 发布工作
// @Description 发布工作
// @Tags wx
// @Param type query string true "工种 0(点工),1(包工) 必填"
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
		work.BasicWork.PublishTime = time.Now().Unix()

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
		work.BasicWork.PublishTime = time.Now().Unix()

		err = tx.Create(&work).Error
		// 数据库错误
		if common.FuncHandler(c, err, nil, 20002) {
			// 发生错误时回滚事务
			tx.Rollback()
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, controller.Message{
			Msg: "发布成功",
		})
	}

}

// SearchWork 搜索工作
// @Summary 搜索工作
// @Description 搜索工作，需要某个筛选参数就加上，否则可以不加；按发布时间降序排序
// @Tags wx
// @Param location query string false "二级位置信息 选填"
// @Param need query string false "所需工种 选填"
// @Param type query string false "工程类别 选填"
// @Param page query string true "页码，从1开始 必填"
// @Param limit query string true "每页记录数 必填"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/work/search [get]
func SearchWork(c *gin.Context) {
	projectType := c.Query("type")
	need := c.Query("need")
	location := c.Query("location")

	page, err := strconv.Atoi(c.Query("page"))
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}
	var limit int
	limit, err = strconv.Atoi(c.Query("limit"))
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}
	if common.FuncHandler(c, page > 0 && limit > 0, true, 30001) {
		return
	}

	db := common.GetMySQL()
	dbsearch := common.GetMySQL()
	if projectType != "" {
		var res model.ProjectType
		err := db.Where("name = ?", projectType).First(&res).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, 30000) {
			return
		}

		dbsearch = dbsearch.Where("project_type = ?", projectType)
	}

	if need != "" {
		var res model.WorkerType
		err := db.Where("name = ?", need).First(&res).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, 30001) {
			return
		}

		dbsearch = dbsearch.Where("worker_type = ?", need)
	}

	if location != "" {
		dbsearch = dbsearch.Where("location = ? ", location)
	}

	var works []model.Work
	err = dbsearch.Find(&works).Error

	count := len(works)
	totalPage := count / limit
	if count%limit != 0 {
		totalPage++
	}

	dbsearch = dbsearch.Order("publish_time desc").Limit(limit).Offset((page - 1) * limit)
	err = dbsearch.Find(&works).Error

	// 找不到数据
	if err != nil {
		c.JSON(http.StatusOK, controller.Message{
			Msg: "无数据",
		})
	} else {
		var mWork []interface{}

		for _, work := range works {
			if work.PricingMode == "点工" {
				var dianWorkRet model.DianWorkReturn
				dianWorkRet.ID = work.ID
				dianWorkRet.BasicWork = work.BasicWork
				dianWorkRet.Treatment = work.Treatment

				locationID := work.LocationID
				dianID := work.Fid

				var locationInfo model.LocationInfo
				err := db.First(&locationInfo, locationID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, 20002) {
					return
				}
				dianWorkRet.LocationInfoWrapper = locationInfo.LocationInfoWrapper

				var dianWork model.DianWork
				err = db.First(&dianWork, dianID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, 20002) {
					return
				}
				dianWorkRet.DianWorkOther = dianWork.DianWorkOther

				mWork = append(mWork, dianWorkRet)
			} else {
				var dianWorkRet model.BaoWorkReturn
				dianWorkRet.ID = work.ID
				dianWorkRet.BasicWork = work.BasicWork
				dianWorkRet.Treatment = work.Treatment

				locationID := work.LocationID
				baoID := work.Fid

				var locationInfo model.LocationInfo
				err := db.First(&locationInfo, locationID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, 20002) {
					return
				}
				dianWorkRet.LocationInfoWrapper = locationInfo.LocationInfoWrapper

				var baoWork model.BaoWork
				err = db.First(&baoWork, baoID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, 20002) {
					return
				}
				dianWorkRet.BaoWorkOther = baoWork.BaoWorkOther

				mWork = append(mWork, dianWorkRet)
			}
		}

		var ret map[string]interface{}
		ret = make(map[string]interface{})
		ret["total_pages"] = totalPage
		ret["current_page"] = page
		ret["data"] = mWork
		c.JSON(http.StatusOK, controller.Message{
			Data: ret,
		})
	}
}
