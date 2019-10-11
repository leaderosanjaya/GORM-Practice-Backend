package user

import (
	"github.com/jinzhu/gorm"
)

// Handler for user handler db
type Handler struct {
	DB *gorm.DB
}

// Del for uid in deleting user
type Del struct {
	UID uint `json:"uid"`
}

// Credential structure
type Credential struct {
	Password string `json:"password"`
}
