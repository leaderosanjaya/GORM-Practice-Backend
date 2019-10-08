package key

import "github.com/jinzhu/gorm"

// Handler handler struct
type Handler struct {
	DB *gorm.DB
}

// Del struct, was KeyDel
type Del struct {
	UID uint `json:"uid"`
}

type Assign struct {
	UID uint `json:"user_id"`
}

// JSONMessage struct
type JSONMessage struct {
	Status    string `json:"status"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
