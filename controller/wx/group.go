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

// NewGroup 创建新班组
// @Summary 创建新班组
// @Description 创建新班组
// @Tags 班组相关
// @Param user body model.GroupRequest true "创建新班组"
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
	newGroup.OwnerID = groupReq.OwnerID
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

// JoinGroup 加入班组
// @Summary 加入班组
// @Description 加入班组
// @Tags 班组相关
// @Param user_id query int true "用户id"
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
	// 检查是否已经在班组
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
		Data: "加入班组成功",
	})
}

// InGroup 查询自己参与的班组
// @Summary 查询自己参与的班组
// @Description 查询自己参与的班组，包括自己创建和加入的
// @Tags 班组相关
// @Param user_id query int true "用户id"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/in_group [get]
func InGroup(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}

	db := common.GetMySQL()
	var groupInfos []model.GroupInfo

	// 查询自己创建的班组
	var groups []model.Group
	err = db.Where("owner_id = ?", userID).Find(&groups).Error

	if err == nil {
		for _, group := range groups {
			var groupInfo model.GroupInfo
			groupInfo.ID = group.ID
			groupInfo.GroupName = group.GroupName
			groupInfo.IsOwner = true
			ownerID := group.OwnerID

			var owner model.WxUser
			err = db.First(&owner, ownerID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, 20002) {
				return
			}

			groupInfo.Owner = owner
			groupInfos = append(groupInfos, groupInfo)
		}
	}
	// 查询自己参与的班组
	var groupMembers []model.GroupMember
	err = db.Where("member_id = ?", userID).Find(&groupMembers).Error

	if err == nil {
		for _, groupMember := range groupMembers {
			groupID := groupMember.GroupID

			var group model.Group
			err = db.First(&group, groupID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, 20002) {
				return
			}

			var groupInfo model.GroupInfo
			groupInfo.ID = group.ID
			groupInfo.GroupName = group.GroupName
			groupInfo.IsOwner = false
			ownerID := group.OwnerID

			var owner model.WxUser
			err = db.First(&owner, ownerID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, 20002) {
				return
			}

			groupInfo.Owner = owner
			groupInfos = append(groupInfos, groupInfo)
		}
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: groupInfos,
	})
}

// GroupMember 获取班组所有成员
// @Summary 获取班组所有成员
// @Description 获取班组所有成员
// @Tags 班组相关
// @Param group_id query int true "班组id"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/group_member [get]
func GroupMember(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}

	db := common.GetMySQL()
	var groupMemberInfos []model.GroupMemberInfo

	// 群主信息
	var group model.Group
	err = db.First(&group, groupID).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, 20002) {
		return
	}

	var user model.WxUser
	err = db.First(&user, group.OwnerID).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, 20002) {
		return
	}
	var groupMemberInfo model.GroupMemberInfo
	groupMemberInfo.IsOwner = true
	groupMemberInfo.WxUser = user
	groupMemberInfos = append(groupMemberInfos, groupMemberInfo)

	// 班组成员信息
	var groupMembers []model.GroupMember
	err = db.Where("group_id = ?", groupID).Find(&groupMembers).Error
	if err == nil {
		for _, groupMember := range groupMembers {
			memberID := groupMember.MemberID

			var user model.WxUser
			err = db.First(&user, memberID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, 20002) {
				return
			}
			var groupMemberInfo model.GroupMemberInfo
			groupMemberInfo.IsOwner = false
			groupMemberInfo.WxUser = user
			groupMemberInfos = append(groupMemberInfos, groupMemberInfo)
		}
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: groupMemberInfos,
	})
}

// DeleteMember 删除班组中某个成员
// @Summary 删除班组中某个成员
// @Description 删除班组中某个成员
// @Tags 班组相关
// @Param group_id query int true "班组id"
// @Param user_id query int true "删除用户id"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/delete_member [get]
func DeleteMember(c *gin.Context) {
	var groupID int64
	var userID int64
	var err error

	groupID, err = strconv.ParseInt(c.Query("group_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}
	userID, err = strconv.ParseInt(c.Query("user_id"), 10, 64)
	if common.FuncHandler(c, err, nil, 20001) {
		return
	}

	db := common.GetMySQL()
	tx := db.Begin()

	err = db.Delete(model.GroupMember{}, "group_id = ? AND member_id = ?", groupID, userID).Error

	if common.FuncHandler(c, err, nil, 20002) {
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Data: "删除成功",
	})
}
