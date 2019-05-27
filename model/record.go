package model

// CommonRecord 通用字段
type CommonRecord struct {
	WorkerID   int64  `json:"worker_id" example:"2"`
	GroupID    int64  `json:"group_id" example:"1"`
	RecordDate string `json:"record_date" example:"2019-05-19"`
	Remark     string `json:"remark"`
}

// Record 工作记录数据库字段
type Record struct {
	ID      int64 `json:"id"`
	AdderID int64 `json:"adder_id"`
	CommonRecord
	RecordType  int64 `json:"record_type"`
	RecordID    int64 `json:"record_id"`
	AddTime     int64 `json:"add_time"`
	IsConfirm   int64 `json:"is_confirm"`
	ConfirmTime int64 `json:"confirm_time"`
}

// RetItemInfo item 记录返回信息
type RetItemInfo struct {
	AdderInfo  WxUserInfo   `json:"adder_info"`
	WorkerInfo WxUserInfo   `json:"worker_info"`
	GroupInfo  GroupRequest `json:"group_info"`
	RecordID   int64        `json:"record_id"`
	RecordDate string       `json:"record_date" example:"2019-05-19"`
	Remark     string       `json:"remark"`
	Subitem    string       `json:"subitem" example:"刷墙"`
	Quantity   float64      `json:"quantity" example:"1"`
	Unit       string       `json:"unit" example:"平方米"`
	AddTime    int64        `json:"add_time"`
	IsConfirm  int64        `json:"is_confirm"`
	Type       int64        `json:"type"`
}

// RetHourInfo hour 记录返回信息
type RetHourInfo struct {
	AdderInfo      WxUserInfo   `json:"adder_info"`
	WorkerInfo     WxUserInfo   `json:"worker_info"`
	GroupInfo      GroupRequest `json:"group_info"`
	RecordID       int64        `json:"record_id"`
	RecordDate     string       `json:"record_date" example:"2019-05-19"`
	Remark         string       `json:"remark"`
	WorkHours      float64      `json:"work_hours" example:"1.5"`
	ExtraWorkHours float64      `json:"extra_work_hours" example:"1"`
	AddTime        int64        `json:"add_time"`
	IsConfirm      int64        `json:"is_confirm"`
	Type           int64        `json:"type"`
}

// HourRecord 工时数据库字段
type HourRecord struct {
	ID             int64   `json:"id"`
	WorkHours      float64 `json:"work_hours" example:"1.5"`
	ExtraWorkHours float64 `json:"extra_work_hours" example:"1"`
}

// HourRecordRequest 工时请求头
type HourRecordRequest struct {
	CommonRecord
	WorkHours      float64 `json:"work_hours" example:"1.5"`
	ExtraWorkHours float64 `json:"extra_work_hours" example:"1"`
}

// ItemRecord 分项记录数据库字段
type ItemRecord struct {
	ID       int64   `json:"id"`
	Subitem  string  `json:"subitem" example:"刷墙"`
	Quantity float64 `json:"quantity" example:"1"`
	Unit     string  `json:"unit" example:"平方米"`
}

// ItemRecordRequest 分项记录请求头
type ItemRecordRequest struct {
	CommonRecord
	Subitem  string  `json:"subitem" example:"刷墙"`
	Quantity float64 `json:"quantity" example:"1"`
	Unit     string  `json:"unit" example:"平方米"`
}
