package model

// Group 群组数据库字段
type Group struct {
	ID      int64 `json:"id"`
	OwnerID int64 `json:"owner_id"`
	GroupRequest
	GroupKey string `json:"group_key"`
}

// GroupRet 详细的班组数据，多了创建者信息
type GroupRet struct {
	ID        int64      `json:"id"`
	GroupName string     `json:"group_name"`
	Owner     WxUserInfo `json:"owner"`
	IsOwner   bool       `json:"is_owner"`
}

// GroupRequest 群组请求字段
type GroupRequest struct {
	GroupName string `json:"group_name"`
}

// GroupMember 群组成员数据库字段
type GroupMember struct {
	ID       int64 `json:"id"`
	GroupID  int64 `json:"group_id"`
	MemberID int64 `json:"member_id"`
}

// GroupMemberRet 群组成员详细信息
type GroupMemberRet struct {
	WxUserInfo
	IsOwner bool `json:"is_owner"`
}
