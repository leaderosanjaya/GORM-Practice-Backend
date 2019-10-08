package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
)

// Token structure object
type Token struct {
	UserID uint
	Name string
	Email string
	Role int
	*jwt.StandardClaims
}

// Handler handler struct
type Handler struct {
	DB *gorm.DB
}