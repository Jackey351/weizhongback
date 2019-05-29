package model

// LocationInfoReq 位置信息请求字段
type LocationInfoReq struct {
	Title     string `json:"title" example:"人民广场"`
	Addr      string `json:"addr" example:"湖北省襄阳市樊城区武商广场店对面人民公园"`
	Latitude  string `json:"latitude" example:"32.04278"`
	Longitude string `json:"longitude" example:"112.15519"`
}

// LocationInfo 位置信息数据库字段
type LocationInfo struct {
	ID int64
	LocationInfoReq
}

// BasicWork 基本信息，不管啥类型都有的字段且出现在请求信息中
type BasicWork struct {
	ConstructionCompany string `json:"construction_company" example:"飞燕工程队"`
	Desc                string `json:"desc" example:"包吃包住"`
	Location            string `json:"location" example:"湖北省襄阳市"`
	WorkerType          string `json:"need" example:"钢筋工"`
	ProjectName         string `json:"project_name" example:"主楼建造"`
	ProjectType         string `json:"type" example:"消防"`
}

// WorkReq 通用工作请求字段
type WorkReq struct {
	BasicWork
	Treatment       []string        `json:"final_treatment"`
	LocationInfoReq LocationInfoReq `json:"location_info"`
}

// Work 与数据表work字段对应
type Work struct {
	ID     int64 `json:"work_id"`
	UserID int64 `json:"user_id"`
	BasicWork
	PricingMode int64  `json:"pricing_mode"`
	PublishTime int64  `json:"publish_time"`
	Treatment   string `json:"treatment"`
	LocationID  int64
	Fid         int64
}

// DianWorkRet 点工返回字段
type DianWorkRet struct {
	ID            int64      `json:"work_id"`
	PublisherInfo WxUserInfo `json:"publisher_info"`
	BasicWork
	PricingMode     int64           `json:"pricing_mode"`
	PublishTime     int64           `json:"publish_time"`
	Treatment       []string        `json:"final_treatment"`
	LocationInfoReq LocationInfoReq `json:"location_info"`
	DianWorkOther
}

// BaoWorkRet 包工返回字段
type BaoWorkRet struct {
	ID            int64      `json:"work_id"`
	PublisherInfo WxUserInfo `json:"publisher_info"`
	BasicWork
	PricingMode     int64           `json:"pricing_mode"`
	PublishTime     int64           `json:"publish_time"`
	Treatment       []string        `json:"final_treatment"`
	LocationInfoReq LocationInfoReq `json:"location_info"`
	BaoWorkOther
}

// TujiWorkRet 突击队返回字段
type TujiWorkRet struct {
	ID            int64      `json:"work_id"`
	PublisherInfo WxUserInfo `json:"publisher_info"`
	BasicWork
	PricingMode     int64           `json:"pricing_mode"`
	PublishTime     int64           `json:"publish_time"`
	Treatment       []string        `json:"final_treatment"`
	LocationInfoReq LocationInfoReq `json:"location_info"`
	TujiWorkOther
}

// DianWorkOther 点工额外信息
type DianWorkOther struct {
	RequireNum string `json:"required_people" example:"11"`
	MaxWage    string `json:"max_wage" example:"200"`
	MinWage    string `json:"min_wage" example:"100"`
	Settlement string `json:"settlement" example:"月薪"`
}

// DianWorkReq 点工请求字段
type DianWorkReq struct {
	WorkReq
	DianWorkOther
}

// DianWork 点工数据库字段
type DianWork struct {
	ID int64 `json:"id"`
	DianWorkOther
}

// BaoWorkOther 包工额外信息
type BaoWorkOther struct {
	Scale      string `json:"scale" example:"大"`
	TotlePrice string `json:"totle_price" example:"300"`
	Unit       string `json:"unit" example:"元/平方米"`
	UnitPrice  string `json:"unit_price" example:"2"`
}

// BaoWorkReq 包工请求字段
type BaoWorkReq struct {
	WorkReq
	BaoWorkOther
}

// BaoWork 包工数据库字段
type BaoWork struct {
	ID int64
	BaoWorkOther
}

// TujiWorkOther 突击队额外信息
type TujiWorkOther struct {
	Num       int64  `json:"required_people" example:"12"`
	StartDate string `json:"work_date" example:"2019-05-12"`
	Days      int64  `json:"work_days" example:"10"`
	Time      string `json:"work_time" example:"8"`
	Money     int64  `json:"money" example:"80"`
}

// TujiWorkReq 突击队
type TujiWorkReq struct {
	WorkReq
	TujiWorkOther
}

// TujiWork 突击队数据库字段
type TujiWork struct {
	ID int64
	TujiWorkOther
}
