package model

// WorkerType 工种类型
type WorkerType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProjectType 工程类别
type ProjectType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
