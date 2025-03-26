package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.StandardClaims
}
