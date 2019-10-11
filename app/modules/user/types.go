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

// JSONMessage structure to define message
type JSONMessage struct {
	Status    string `json:"status,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
}

// Credential structure
type Credential struct {
	Password string `json:"password"`
}
