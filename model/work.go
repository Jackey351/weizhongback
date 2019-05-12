package model

// LocationInfoWrapper 位置信息wrapper
type LocationInfoWrapper struct {
	Title     string `json:"title" example:"人民广场"`
	Addr      string `json:"addr" example:"湖北省襄阳市樊城区武商广场店对面人民公园"`
	Latitude  string `json:"latitude" example:"32.04278"`
	Longitude string `json:"longitude" example:"112.15519"`
}

// LocationInfo 位置信息wrapper
type LocationInfo struct {
	ID int64
	LocationInfoWrapper
}

// BasicWork 基本信息，不管啥类型都有的字段且出现在请求信息中
type BasicWork struct {
	ConstructionCompany string `json:"construction_company" example:"飞燕工程队"`
	Desc                string `json:"desc" example:"包吃包住"`
	Location            string `json:"location" example:"湖北省襄阳市"`
	WorkerType          string `json:"need" example:"钢筋工"`
	ProjectName         string `json:"project_name" example:"主楼建造"`
	ProjectType         string `json:"type" example:"消防"`
	UserID              int64  `json:"user_id" example:"1"`
	Nav                 string `json:"nav" example:"工人"`
}

// WorkWrapper 工作请求wrapper
type WorkWrapper struct {
	BasicWork
	Treatment           []string            `json:"final_treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
}

// Work 与数据表work字段对应
type Work struct {
	ID int64 `json:"work_id"`
	BasicWork
	PricingMode string `json:"pricing_mode"`
	PublishTime int64  `json:"publish_time"`
	Treatment   string `json:"treatment"`
	LocationID  int64
	Fid         int64
}

// DianWorkReturn 点工作为返回值
type DianWorkReturn struct {
	ID int64 `json:"work_id"`
	BasicWork
	PricingMode         string              `json:"pricing_mode"`
	PublishTime         int64               `json:"publish_time"`
	Treatment           []string            `json:"final_treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
	DianWorkOther
}

// BaoWorkReturn 包工作为返回值
type BaoWorkReturn struct {
	ID int64 `json:"work_id"`
	BasicWork
	PricingMode         string              `json:"pricing_mode"`
	PublishTime         int64               `json:"publish_time"`
	Treatment           []string            `json:"final_treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
	BaoWorkOther
}

// TujiWorkReturn 突击队作为返回值
type TujiWorkReturn struct {
	ID int64 `json:"work_id"`
	BasicWork
	PricingMode         string              `json:"pricing_mode"`
	PublishTime         int64               `json:"publish_time"`
	Treatment           []string            `json:"final_treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
	TujiWorkOther
}

// DianWorkOther 点工额外信息
type DianWorkOther struct {
	RequireNum string `json:"required_people" example:"11"`
	MaxWage    string `json:"max_wage" example:"200"`
	MinWage    string `json:"min_wage" example:"100"`
	Settlement string `json:"settlement" example:"月薪"`
}

// DianWorkWrapper 点工
type DianWorkWrapper struct {
	WorkWrapper
	DianWorkOther
}

// DianWork 点工有id
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

// BaoWorkWrapper 包工
type BaoWorkWrapper struct {
	WorkWrapper
	BaoWorkOther
}

// BaoWork 包工有id
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

// TujiWorkWrapper 突击队
type TujiWorkWrapper struct {
	WorkWrapper
	TujiWorkOther
}

// TujiWork 突击队有id
type TujiWork struct {
	ID int64
	TujiWorkOther
}
