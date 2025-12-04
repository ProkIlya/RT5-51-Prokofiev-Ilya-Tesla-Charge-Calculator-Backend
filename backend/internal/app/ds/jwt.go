package ds

import (
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	UserID      uint   `json:"user_id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
	jwt.StandardClaims
}
