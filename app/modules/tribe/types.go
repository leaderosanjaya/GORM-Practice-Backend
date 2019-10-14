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

// Assign user_id
type Assign struct {
	UID uint `json:"user_id"`
	PlatformID uint `json:"platform"`
}
