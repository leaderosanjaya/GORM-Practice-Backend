package models

import "time"

//User struct, refer to DB data
type User struct {
	ID         uint          `json:"user_id" gorm:"primary_key;column:user_id"`
	CreatedAt  time.Time     `json:"-"`
	UpdatedAt  time.Time     `json:"-"`
	FirstName  string        `json:"first_name" gorm:"type:varchar(20);not null"`
	LastName   string        `json:"last_name" gorm:"type:varchar(20);not null"`
	Email      string        `json:"email" gorm:"type:varchar(50);unique;not null"`
	Password   string        `json:"-" gorm:"type:varchar(255);not null"`
	Role       int           `json:"role" gorm:"default:0"`
	Keys       []Key         `json:"keys" gorm:"foreignkey:UserID"`
	Tribes     []TribeAssign `json:"tribes"`
	SharedKeys []KeyShares   `json:"shared_keys"`
}

// Tribe struct, tribe data model in DB
type Tribe struct {
	ID          uint          `json:"tribe_id" gorm:"primary_key;column:tribe_id"`
	CreatedAt   time.Time     `json:"-"`
	UpdatedAt   time.Time     `json:"-"`
	TribeName   string        `json:"tribe_name" gorm:"type:varchar(50);not null;unique"`
	LeadID      uint          `json:"lead_id"` //use lead_id as foreign key
	Description string        `json:"description" gorm:"type:varchar(200)"`
	TotalMember int           `json:"total_member" gorm:"not null;default:1"`
	TotalKey    int           `json:"total_key" gorm:"not null;default:0"`
	Keys        []Key         `json:"keys" gorm:"foreignkey:TribeID`
	Members     []TribeAssign `json:"members"`
}

// Key struct, key data model in DB
type Key struct {
	ID          uint        `json:"key_id" gorm:"primary_key;column:key_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"-"`
	KeyName     string      `json:"key_name" gorm:"type:varchar(50);not null"`
	KeyValue    string      `json:"key_value" gorm:"type:varchar(300);not null"`
	KeyType     string      `json:"key_type" gorm:"type:varchar(15);not null;default:'STRING'"`
	Description string      `json:"description" gorm:"type:varchar(200);not null"`
	Platform    string      `json:"platform" gorm:"type:varchar(50);not null"`
	ExpireDate  time.Time   `json:"expire_date" gorm:"not null"`
	User        User        `json:"-"`
	UserID      uint        `json:"user_id" gorm:"not null"`
	Tribe       Tribe       `json:"-"`
	TribeID     uint        `json:"tribe_id" gorm:"not null"`
	AppVersion  string      `json:"app_version" gorm:"type:varchar(20);not null"`
	Status      string      `json:"status" gorm:"type:varchar(20);not null"`
	Shares      []KeyShares `json:"shares"`
}

type KeyShares struct {
	UserID uint
	KeyID  uint
}

type TribeAssign struct {
	UserID  uint
	TribeID uint
}
