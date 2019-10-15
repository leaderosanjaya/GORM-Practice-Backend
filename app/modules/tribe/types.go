package tribe

import "github.com/jinzhu/gorm"

// Handler struct object
type Handler struct {
	DB *gorm.DB
}

// Del Struct object
type Del struct {
	UID uint `json:"uid"`
}

// Assign user_id
type Assign struct {
	UID uint `json:"user_id"`
}

// TribeCreate struct object
type TribeCreate struct {
	TribeName   string `json:"tribe_name"`
	LeadID      uint   `json:"lead_id"`
	Description string `json:"description"`
}
