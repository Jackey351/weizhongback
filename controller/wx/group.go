package wx

import (
	"math/rand"
	"net/http"
	"time"
	"yanfei_backend/common"
	"yanfei_backend/controller"
	"yanfei_backend/model"

	"github.com/gin-gonic/gin"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// NewGroupKey 生成随机的不重复的群组id
func NewGroupKey() string {
	n := 4
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	newGroupKey := string(b)
	db := common.GetMySQL()

	var existGroup model.Group
	db.Where("group_key = ?", newGroupKey).First(&existGroup)

	if existGroup.ID != 0 {
		return NewGroupKey()
	}
	return newGroupKey
}

// NewGroup 创建新群组
// @Summary 创建新群组
// @Description 创建新群组
// @Tags wx
// @Param user body model.GroupRequest true "create a new group"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/new_group [post]
func NewGroup(c *gin.Context) {
	var groupReq model.GroupRequest

	// 获取数据失败
	if common.FuncHandler(c, c.BindJSON(&groupReq), nil, 20001) {
		return
	}

	var newGroup model.Group
	newGroup.GroupName = groupReq.GroupName
	newGroup.UID = groupReq.UID
	newGroup.GroupKey = NewGroupKey()

	db := common.GetMySQL()
	tx := db.Begin()

	err := tx.Create(&newGroup).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, controller.Message{
		Data: newGroup,
	})
}
