package key

import "github.com/jinzhu/gorm"

type Handler struct {
	DB *gorm.DB
}

type KeyDel struct {
	UID uint `json:"uid"`
}

type JSONMessage struct {
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
