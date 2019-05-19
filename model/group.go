package model

// Group 群组数据库字段
type Group struct {
	ID int64 `json:"id"`
	GroupRequest
	GroupKey string `json:"group_key"`
}

// GroupInfo 详细的班组数据，多了创建者信息
type GroupInfo struct {
	ID        int64  `json:"id"`
	GroupName string `json:"group_name"`
	Owner     WxUser `json:"owner"`
	IsOwner   bool   `json:"is_owner"`
}

// GroupRequest 群组请求字段
type GroupRequest struct {
	OwnerID   int64  `json:"owner_id"`
	GroupName string `json:"group_name"`
}

// GroupMember 群组成员数据库字段
type GroupMember struct {
	ID       int64 `json:"id"`
	GroupID  int64 `json:"group_id"`
	MemberID int64 `json:"member_id"`
}
