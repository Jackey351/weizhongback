package model

// User user info
type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
