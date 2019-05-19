package model

// WorkType 工种类型
type WorkType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProjectType 工程类别
type ProjectType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
