package models

import "time"

//User struct, refer to DB data
type User struct {
	ID        uint      `json:"user_id" gorm:"primary_key;column:user_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	FirstName string    `json:"first_name" gorm:"type:varchar(20);not null"`
	LastName  string    `json:"last_name" gorm:"type:varchar(20);not null"`
	Email     string    `json:"email" gorm:"type:varchar(50);unique;not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	Role      int       `json:"role" gorm:"default:0"`
	// Platform   int           `json:"platform" gorm:"default:0"`
	Keys       []Key         `json:"keys" gorm:"foreignkey:UserID"`
	Tribes     []TribeAssign `json:"tribes"`
	SharedKeys []KeyShares   `json:"shared_keys"`
}

// Tribe struct, tribe data model in DB
type Tribe struct {
	ID          uint              `json:"tribe_id" gorm:"primary_key;column:tribe_id"`
	CreatedAt   time.Time         `json:"-"`
	UpdatedAt   time.Time         `json:"-"`
	TribeName   string            `json:"tribe_name" gorm:"type:varchar(255);not null;unique"`
	Leads       []TribeLeadAssign `json:"tribe_leads"` //use lead_id as foreign key
	Description string            `json:"description" gorm:"type:text"`
	TotalMember int               `json:"total_member" gorm:"not null;default:0"`
	TotalKey    int               `json:"total_key" gorm:"not null;default:0"`
	Keys        []Key             `json:"keys" gorm:"foreignkey:TribeID"`
	Members     []TribeAssign     `json:"members"`
}

// Key struct, key data model in DB
type Key struct {
	ID          uint        `json:"key_id" gorm:"primary_key;column:key_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	KeyName     string      `json:"key_name" gorm:"type:varchar(255);not null;unique"`
	KeyValue    string      `json:"key_value" gorm:"type:text;not null"`
	KeyType     string      `json:"key_type" gorm:"type:varchar(20);not null;default:'STRING'"`
	Description string      `json:"description" gorm:"type:text;not null"`
	Platform    string      `json:"platform" gorm:"type:varchar(50);not null"`
	ExpireDate  time.Time   `json:"expire_date" gorm:"not null"`
	UserID      uint        `json:"user_id" gorm:"not null"`
	TribeID     uint        `json:"tribe_id" gorm:"not null"`
	AppVersion  string      `json:"app_version" gorm:"type:varchar(20);not null"`
	Status      string      `json:"status" gorm:"type:varchar(20);not null"`
	Shares      []KeyShares `json:"shares"`
}

// KeyShares user association with key
type KeyShares struct {
	UserID uint `gorm:"primary_key"`
	KeyID  uint `gorm:"primary_key"`
}

// TribeAssign user association with tribe
type TribeAssign struct {
	UserID  uint `gorm:"primary_key"`
	TribeID uint `gorm:"primary_key"`
	Platform uint `gorm:"not null; default:0"`
}

// TribeLeadAssign user lead
type TribeLeadAssign struct {
	LeadID  uint `gorm:"primary_key"`
	TribeID uint `gorm:"primary_key"`
}
