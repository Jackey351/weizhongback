package wx

import (
	"math/rand"
	"net/http"
	"strconv"
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

// JoinGroup 加入群组
// @Summary 加入群组
// @Description 加入群组
// @Tags wx
// @Param user_id query string true "用户id"
// @Param group_key query string true "群组入群口令"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/join_group [get]
func JoinGroup(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}
	groupKey := c.Query("group_key")

	db := common.GetMySQL()
	// 检查userID是否存在
	var existUser model.WxUser
	db.First(&existUser, userID)
	if common.FuncHandler(c, existUser.ID != 0, true, 40000) {
		return
	}

	// 检查groupKey是否存在
	var existGroup model.Group
	db.Where("group_key = ?", groupKey).First(&existGroup)
	if common.FuncHandler(c, existGroup.ID != 0, true, 50000) {
		return
	}

	groupID := existGroup.ID
	// 检查是否已经在群组
	var existGroupMember model.GroupMember
	db.Where("group_id = ? AND member_id = ?", groupID, userID).First(&existGroupMember)
	if common.FuncHandler(c, existGroupMember.ID == 0, true, 50001) {
		return
	}

	var newGroupMember model.GroupMember
	newGroupMember.MemberID = userID
	newGroupMember.GroupID = groupID

	tx := db.Begin()

	err = tx.Create(&newGroupMember).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, 20002) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, controller.Message{
		Data: "加入群组成功",
	})
}
