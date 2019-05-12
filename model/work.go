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

// BasicWork 基本信息
type BasicWork struct {
	ConstructionCompany string `json:"construction_company" example:"飞燕工程队"`
	Desc                string `json:"desc" example:"包吃包住"`
	Location            string `json:"location" example:"湖北省襄阳市"`
	WorkerType          string `json:"need" example:"钢筋工"`
	PricingMode         string `json:"pricing_mode" example:"点工"`
	ProjectName         string `json:"project_name" example:"主楼建造"`
	ProjectType         string `json:"type" example:"消防"`
	UserID              int64  `json:"user_id" example:"1"`
	PublishTime         int64  `json:"publish_time"`
}

// WorkWrapper 工作请求wrapper
type WorkWrapper struct {
	BasicWork
	Treatment           []string            `json:"final_treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
}

// Work 有id
type Work struct {
	ID int64 `json:"work_id"`
	BasicWork
	Treatment  string `json:"treatment"`
	LocationID int64
	Fid        int64
}

// DianWorkReturn 点工作为返回值
type DianWorkReturn struct {
	ID int64 `json:"work_id"`
	BasicWork
	Treatment           string              `json:"treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
	DianWorkOther
}

// BaoWorkReturn 包工作为返回值
type BaoWorkReturn struct {
	ID int64 `json:"work_id"`
	BasicWork
	Treatment           string              `json:"treatment"`
	LocationInfoWrapper LocationInfoWrapper `json:"location_info"`
	BaoWorkOther
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
