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

// GroupExistByID 根据班组id判断班组是否存在，不存在直接返回GroupNoExist
func GroupExistByID(c *gin.Context, groupID int64) interface{} {
	db := common.GetMySQL()

	var existGroup model.Group
	err := db.First(&existGroup, groupID).Error
	if common.FuncHandler(c, err, nil, common.GroupNoExist) {
		return false
	}

	return existGroup
}

// GroupExistByKey 根据班组key判断班组是否存在，不存在直接返回GroupNoExist
func GroupExistByKey(c *gin.Context, groupKey string) interface{} {
	db := common.GetMySQL()

	var existGroup model.Group
	err := db.Where("group_key = ?", groupKey).First(&existGroup).Error
	if common.FuncHandler(c, err, nil, common.GroupNoExist) {
		return false
	}

	return existGroup
}

// NewGroup 创建新班组
// @Summary 创建新班组
// @Description 创建新班组
// @Tags 班组相关
// @Param token header string true "token"
// @Param user body model.GroupRequest true "创建新班组"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/new_group [post]
func NewGroup(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	var groupReq model.GroupRequest

	// 获取数据失败
	if common.FuncHandler(c, c.BindJSON(&groupReq), nil, common.ParameterError) {
		return
	}

	if _, ok := UserExist(c, userID).(model.WxUser); !ok {
		return
	}

	var newGroup model.Group
	newGroup.GroupRequest = groupReq
	newGroup.OwnerID = userID
	newGroup.GroupKey = NewGroupKey()

	db := common.GetMySQL()
	tx := db.Begin()

	err := tx.Create(&newGroup).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
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
// @Param token header string true "token"
// @Param group_key query string true "群组入群口令"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/join_group [get]
func JoinGroup(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	groupKey := c.Query("group_key")

	// 检查userID是否存在
	if _, ok := UserExist(c, userID).(model.WxUser); !ok {
		return
	}

	// 检查groupKey是否存在
	var existGroup model.Group
	var ok bool
	if existGroup, ok = GroupExistByKey(c, groupKey).(model.Group); !ok {
		return
	}

	if common.FuncHandler(c, existGroup.OwnerID != userID, true, common.HasInGroup) {
		return
	}

	db := common.GetMySQL()

	groupID := existGroup.ID
	// 检查是否已经在班组
	var existGroupMember model.GroupMember
	err := db.Where("group_id = ? AND member_id = ?", groupID, userID).First(&existGroupMember).Error
	if common.FuncHandler(c, err != nil, true, common.HasInGroup) {
		return
	}

	var newGroupMember model.GroupMember
	newGroupMember.MemberID = userID
	newGroupMember.GroupID = groupID

	tx := db.Begin()

	err = tx.Create(&newGroupMember).Error
	// 数据库错误
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		// 发生错误时回滚事务
		tx.Rollback()
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, controller.Message{
		Msg: "加入班组成功",
	})
}

// InGroup 查询自己参与的班组
// @Summary 查询自己参与的班组
// @Description 查询自己参与的班组，包括自己创建和加入的
// @Tags 班组相关
// @Param token header string true "token"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/in_group [get]
func InGroup(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	userID := claims.(*model.CustomClaims).UserID

	// 检查userID是否存在
	if _, ok := UserExist(c, userID).(model.WxUser); !ok {
		return
	}

	db := common.GetMySQL()
	var groupRets []model.GroupRet

	// 查询自己创建的班组
	var groups []model.Group
	err := db.Where("owner_id = ?", userID).Find(&groups).Error

	if err == nil {
		for _, group := range groups {
			var groupRet model.GroupRet
			groupRet.ID = group.ID
			groupRet.GroupName = group.GroupName
			groupRet.IsOwner = true
			ownerID := group.OwnerID

			var owner model.WxUser
			err = db.First(&owner, ownerID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, common.DatabaseError) {
				return
			}

			groupRet.Owner = owner.WxUserInfo
			groupRets = append(groupRets, groupRet)
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
			if common.FuncHandler(c, err, nil, common.DatabaseError) {
				return
			}

			var groupRet model.GroupRet
			groupRet.ID = group.ID
			groupRet.GroupName = group.GroupName
			groupRet.IsOwner = false
			ownerID := group.OwnerID

			var owner model.WxUser
			err = db.First(&owner, ownerID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, common.DatabaseError) {
				return
			}

			groupRet.Owner = owner.WxUserInfo
			groupRets = append(groupRets, groupRet)
		}
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: groupRets,
	})
}

// GroupMember 获取班组所有成员
// @Summary 获取班组所有成员
// @Description 获取班组所有成员
// @Tags 班组相关
// @Param token header string true "token"
// @Param group_id query int true "班组id"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/group_member [get]
func GroupMember(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if common.FuncHandler(c, err, nil, common.ParameterError) {
		return
	}

	var group model.Group
	var ok bool
	if group, ok = GroupExistByID(c, groupID).(model.Group); !ok {
		return
	}

	db := common.GetMySQL()
	var groupMemberRets []model.GroupMemberRet

	var user model.WxUser
	err = db.First(&user, group.OwnerID).Error
	// 找不到数据
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		return
	}
	var groupMemberRet model.GroupMemberRet
	groupMemberRet.IsOwner = true
	groupMemberRet.WxUserInfo = user.WxUserInfo
	groupMemberRets = append(groupMemberRets, groupMemberRet)

	// 班组成员信息
	var groupMembers []model.GroupMember
	err = db.Where("group_id = ?", groupID).Find(&groupMembers).Error
	if err == nil {
		for _, groupMember := range groupMembers {
			memberID := groupMember.MemberID

			var user model.WxUser
			err = db.First(&user, memberID).Error
			// 找不到数据
			if common.FuncHandler(c, err, nil, common.DatabaseError) {
				return
			}
			var groupMemberRet model.GroupMemberRet
			groupMemberRet.IsOwner = false
			groupMemberRet.WxUserInfo = user.WxUserInfo
			groupMemberRets = append(groupMemberRets, groupMemberRet)
		}
	}

	c.JSON(http.StatusOK, controller.Message{
		Data: groupMemberRets,
	})
}

// DeleteMember 删除班组中某个成员
// @Summary 删除班组中某个成员
// @Description 删除班组中某个成员
// @Tags 班组相关
// @Param token header string true "token"
// @Param group_id query int true "班组id"
// @Param user_id query int true "删除用户id"
// @Accept json
// @Produce json
// @Success 200 {object} controller.Message
// @Router /wx/group/delete_member [get]
func DeleteMember(c *gin.Context) {
	claims, exist := c.Get("claims")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}
	reqUserID := claims.(*model.CustomClaims).UserID

	var groupID int64
	var userID int64
	var err error

	groupID, err = strconv.ParseInt(c.Query("group_id"), 10, 64)
	if common.FuncHandler(c, err, nil, common.ParameterError) {
		return
	}
	userID, err = strconv.ParseInt(c.Query("user_id"), 10, 64)
	if common.FuncHandler(c, err, nil, common.ParameterError) {
		return
	}

	if group, ok := GroupExistByID(c, groupID).(model.Group); !ok {
		return
	} else {
		// 检查是否为班组长
		if common.FuncHandler(c, group.OwnerID == reqUserID, true, common.NoPermission) {
			return
		}
	}

	if _, ok := UserExist(c, userID).(model.WxUser); !ok {
		return
	}

	db := common.GetMySQL()
	tx := db.Begin()

	err = db.Delete(model.GroupMember{}, "group_id = ? AND member_id = ?", groupID, userID).Error

	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, controller.Message{
		Msg: "删除成功",
	})
}
