package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Token structure object
type Token struct {
	UserID uint
	Name   string
	Email  string
	Role   int
	*jwt.StandardClaims
}

// Handler handler struct
type Handler struct {
	DB *gorm.DB
}

// Credential object struct
type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims struct
type Claims struct {
	jwt.StandardClaims
}

// JSONMessage struct object
type JSONMessage struct {
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}