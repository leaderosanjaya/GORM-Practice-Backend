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
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
