package model

import "github.com/dgrijalva/jwt-go"

// CustomClaims 载荷
type CustomClaims struct {
	UserID int64 `json:"user_Id"`
	jwt.StandardClaims
}
