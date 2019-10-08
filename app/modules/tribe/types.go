package tribe

import "github.com/jinzhu/gorm"

// Handler struct objet
type Handler struct {
	DB *gorm.DB
}

// Del Struct objet
type Del struct {
	UID uint `json:"uid"`
}

// JSONMessage struct object
type JSONMessage struct {
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
