package tribe

import "github.com/jinzhu/gorm"

type Handler struct {
	DB *gorm.DB
}

type TribeDel struct {
	UID uint `json:"uid"`
}

type JSONMessage struct {
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
