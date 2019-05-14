package model

// Group 群组数据库字段
type Group struct {
	ID        int64  `json:"id"`
	GroupName string `json:"group_name"`
	UID       int64  `json:"uid"`
	GroupKey  string `json:"group_key"`
}

// GroupRequest 群组请求字段
type GroupRequest struct {
	UID       int64  `json:"uid"`
	GroupName string `json:"group_name"`
}
