package wx

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"hackthoon/common"
	"hackthoon/controller"
	"hackthoon/model"
	"hackthoon/storage"

	"github.com/gin-gonic/gin"
)

// 工作计价类型
const (
	DianWork = 0
	BaoWork  = 1
	TujiWork = 2
)

// PublishDianWork 发布点工工作
// @Summary 发布点工工作
// @Description 发布点工工作
// @Tags 工作相关
// @Param token header string true "token"
// @Param dian_work body model.DianWorkReq true "点工招聘数据"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/work/publish_dian [post]
func PublishDianWork(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var dianWorkReq model.DianWorkReq
	if common.FuncHandler(c, c.BindJSON(&dianWorkReq), nil, common.ParameterError) {
		return
	}

	db := common.GetMySQL()

	// 检查userID是否存在
	_, err := storage.UserExist(userID)
	if common.FuncHandler(c, err, nil, common.UserNoExist) {
		return
	}

	// 检查工种和工程类别是否正确
	projectType := dianWorkReq.BasicWork.ProjectType
	workerType := dianWorkReq.BasicWork.WorkerType

	var res model.ProjectType
	err = db.Where("name = ?", projectType).First(&res).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.ProjectTypeNoExist) {
		return
	}

	var res2 model.WorkType
	err = db.Where("name = ?", workerType).First(&res2).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.WorkTypeNoExist) {
		return
	}

	var dianWork model.DianWork
	dianWork.DianWorkOther = dianWorkReq.DianWorkOther

	tx := db.Begin()

	err = tx.Create(&dianWork).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var locationInfo model.LocationInfo
	locationInfo.LocationInfoReq = dianWorkReq.WorkReq.LocationInfoReq
	err = tx.Create(&locationInfo).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var work model.Work
	work.UserID = userID
	work.BasicWork = dianWorkReq.WorkReq.BasicWork
	work.LocationID = locationInfo.ID
	work.Treatment = strings.Join(dianWorkReq.WorkReq.Treatment, ",")
	work.Fid = dianWork.ID
	work.PricingMode = DianWork
	work.PublishTime = time.Now().Unix()

	err = tx.Create(&work).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Msg: "发布成功",
	})
}

// PublishBaoWork 发布包工工作
// @Summary 发布包工工作
// @Description 发布包工工作
// @Tags 工作相关
// @Param token header string true "token"
// @Param bao_work body model.BaoWorkReq true "包工招聘数据"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/work/publish_bao [post]
func PublishBaoWork(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var baoWorkReq model.BaoWorkReq
	if common.FuncHandler(c, c.BindJSON(&baoWorkReq), nil, common.ParameterError) {
		return
	}

	db := common.GetMySQL()

	// 检查userID是否存在
	_, err := storage.UserExist(userID)
	if common.FuncHandler(c, err, nil, common.UserNoExist) {
		return
	}

	// 检查工种和工程类别是否正确
	projectType := baoWorkReq.BasicWork.ProjectType
	workerType := baoWorkReq.BasicWork.WorkerType

	var res model.ProjectType
	err = db.Where("name = ?", projectType).First(&res).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.ProjectTypeNoExist) {
		return
	}

	var res2 model.WorkType
	err = db.Where("name = ?", workerType).First(&res2).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.WorkTypeNoExist) {
		return
	}

	var baoWork model.BaoWork
	baoWork.BaoWorkOther = baoWorkReq.BaoWorkOther

	tx := db.Begin()

	err = tx.Create(&baoWork).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var locationInfo model.LocationInfo
	locationInfo.LocationInfoReq = baoWorkReq.WorkReq.LocationInfoReq
	err = tx.Create(&locationInfo).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var work model.Work
	work.UserID = userID
	work.BasicWork = baoWorkReq.WorkReq.BasicWork
	work.LocationID = locationInfo.ID
	work.Treatment = strings.Join(baoWorkReq.WorkReq.Treatment, ",")
	work.Fid = baoWork.ID
	work.PricingMode = BaoWork
	work.PublishTime = time.Now().Unix()

	err = tx.Create(&work).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Msg: "发布成功",
	})

}

// PublishTujiWork 发布突击队工作
// @Summary 发布突击队工作
// @Description 发布突击队工作
// @Tags 工作相关
// @Param token header string true "token"
// @Param tuji_work body model.TujiWorkReq true "突击队招聘数据"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/work/publish_tuji [post]
func PublishTujiWork(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var tujiWorkReq model.TujiWorkReq
	if common.FuncHandler(c, c.BindJSON(&tujiWorkReq), nil, common.ParameterError) {
		return
	}

	db := common.GetMySQL()

	// 检查userID是否存在
	_, err := storage.UserExist(userID)
	if common.FuncHandler(c, err, nil, common.UserNoExist) {
		return
	}

	// 检查工种和工程类别是否正确
	projectType := tujiWorkReq.BasicWork.ProjectType
	workerType := tujiWorkReq.BasicWork.WorkerType

	var res model.ProjectType
	err = db.Where("name = ?", projectType).First(&res).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.ProjectTypeNoExist) {
		return
	}

	var res2 model.WorkType
	err = db.Where("name = ?", workerType).First(&res2).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.WorkTypeNoExist) {
		return
	}

	var tujiWork model.TujiWork
	tujiWork.TujiWorkOther = tujiWorkReq.TujiWorkOther

	tx := db.Begin()

	err = tx.Create(&tujiWork).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var locationInfo model.LocationInfo
	locationInfo.LocationInfoReq = tujiWorkReq.WorkReq.LocationInfoReq
	err = tx.Create(&locationInfo).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	var work model.Work
	work.BasicWork = tujiWorkReq.WorkReq.BasicWork
	work.UserID = userID
	work.LocationID = locationInfo.ID
	work.Treatment = strings.Join(tujiWorkReq.WorkReq.Treatment, ",")
	work.Fid = tujiWork.ID
	work.PricingMode = TujiWork
	work.PublishTime = time.Now().Unix()

	err = tx.Create(&work).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Msg: "发布成功",
	})
}

// SearchWork 搜索工作
// @Summary 搜索工作
// @Description 搜索工作，需要某个筛选参数就加上，否则可以不加；按发布时间降序排序
// @Tags 工作相关
// @Param token header string true "token"
// @Param location query string false "二级位置信息 选填"
// @Param need query string false "所需工种 选填"
// @Param type query string false "工程类别 选填"
// @Param work_type query int false "工作类别 选填0只返回点工和包工，1只返回突击队，默认为0"
// @Param page query int true "页码，从1开始 必填"
// @Param limit query int true "每页记录数 必填"
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/work/search [get]
func SearchWork(c *gin.Context) {
	projectType := c.Query("type")
	need := c.Query("need")
	location := c.Query("location")

	workTypeQuery := c.Query("work_type")
	var workType int
	if workTypeQuery != "" {
		var err error
		workType, err = strconv.Atoi(workTypeQuery)
		if common.FuncHandler(c, err, nil, common.ParameterError) {
			return
		}
	} else {
		workType = 0
	}

	page, err := strconv.Atoi(c.Query("page"))
	if common.FuncHandler(c, err, nil, common.ParameterError) {
		return
	}
	var limit int
	limit, err = strconv.Atoi(c.Query("limit"))
	if common.FuncHandler(c, err, nil, common.ParameterError) {
		return
	}
	if common.FuncHandler(c, page > 0 && limit > 0, true, common.ParameterError) {
		return
	}

	db := common.GetMySQL()
	dbsearch := common.GetMySQL()
	if projectType != "" {
		var res model.ProjectType
		err := db.Where("name = ?", projectType).First(&res).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, common.ProjectTypeNoExist) {
			return
		}

		dbsearch = dbsearch.Where("project_type = ?", projectType)
	}

	if need != "" {
		var res model.WorkType
		err := db.Where("name = ?", need).First(&res).Error
		// 找不到数据
		if common.FuncHandler(c, err, nil, common.WorkTypeNoExist) {
			return
		}

		dbsearch = dbsearch.Where("worker_type = ?", need)
	}

	if location != "" {
		dbsearch = dbsearch.Where("location = ?", location)
	}

	if workType == 1 {
		dbsearch = dbsearch.Where("pricing_mode = ?", TujiWork)
	} else {
		dbsearch = dbsearch.Where("pricing_mode = ? OR pricing_mode = ?", DianWork, BaoWork)
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
			switch work.PricingMode {
			case DianWork:
				var dianWorkRet model.DianWorkRet
				dianWorkRet.ID = work.ID
				dianWorkRet.BasicWork = work.BasicWork
				dianWorkRet.Treatment = strings.Split(work.Treatment, ",")
				dianWorkRet.PricingMode = work.PricingMode
				dianWorkRet.PublishTime = work.PublishTime

				user, err := storage.UserExist(work.UserID)
				if common.FuncHandler(c, err, nil, common.UserNoExist) {
					return
				}
				dianWorkRet.PublisherInfo = user.WxUserInfo

				locationID := work.LocationID
				dianID := work.Fid

				var locationInfo model.LocationInfo
				err = db.First(&locationInfo, locationID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, common.DatabaseError) {
					return
				}
				dianWorkRet.LocationInfoReq = locationInfo.LocationInfoReq

				var dianWork model.DianWork
				err = db.First(&dianWork, dianID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, common.DatabaseError) {
					return
				}
				dianWorkRet.DianWorkOther = dianWork.DianWorkOther

				mWork = append(mWork, dianWorkRet)

				break

			case BaoWork:
				var baoWorkRet model.BaoWorkRet
				baoWorkRet.ID = work.ID
				baoWorkRet.BasicWork = work.BasicWork
				baoWorkRet.Treatment = strings.Split(work.Treatment, ",")
				baoWorkRet.PricingMode = work.PricingMode
				baoWorkRet.PublishTime = work.PublishTime

				user, err := storage.UserExist(work.UserID)
				if common.FuncHandler(c, err, nil, common.UserNoExist) {
					return
				}
				baoWorkRet.PublisherInfo = user.WxUserInfo

				locationID := work.LocationID
				baoID := work.Fid

				var locationInfo model.LocationInfo
				err = db.First(&locationInfo, locationID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, common.DatabaseError) {
					return
				}
				baoWorkRet.LocationInfoReq = locationInfo.LocationInfoReq

				var baoWork model.BaoWork
				err = db.First(&baoWork, baoID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, common.DatabaseError) {
					return
				}
				baoWorkRet.BaoWorkOther = baoWork.BaoWorkOther

				mWork = append(mWork, baoWorkRet)
				break
			case TujiWork:
				var tujiWorkRet model.TujiWorkRet
				tujiWorkRet.ID = work.ID
				tujiWorkRet.BasicWork = work.BasicWork
				tujiWorkRet.Treatment = strings.Split(work.Treatment, ",")
				tujiWorkRet.PricingMode = work.PricingMode
				tujiWorkRet.PublishTime = work.PublishTime
				locationID := work.LocationID
				tujiID := work.Fid

				user, err := storage.UserExist(work.UserID)
				if common.FuncHandler(c, err, nil, common.UserNoExist) {
					return
				}
				tujiWorkRet.PublisherInfo = user.WxUserInfo

				var locationInfo model.LocationInfo
				err = db.First(&locationInfo, locationID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, common.DatabaseError) {
					return
				}
				tujiWorkRet.LocationInfoReq = locationInfo.LocationInfoReq

				var tujiWork model.TujiWork
				err = db.First(&tujiWork, tujiID).Error
				// 找不到数据
				if common.FuncHandler(c, err, nil, common.DatabaseError) {
					return
				}
				tujiWorkRet.TujiWorkOther = tujiWork.TujiWorkOther

				mWork = append(mWork, tujiWorkRet)
				break
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
