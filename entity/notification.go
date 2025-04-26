package entity

import (
	"gorm.io/gorm"
	"time"
)

type NotificationType string

const (
	NotificationEmail NotificationType = "email"
)

type Notification struct {
	ID        string           `json:"id" gorm:"type:varchar(255);primarykey"`
	Type      NotificationType `json:"type" gorm:"type:varchar(255)" validate:"required"`
	Email     EncryptedField   `json:"email" gorm:"type:varbinary(512)" validate:"required"`
	Title     string           `json:"title" gorm:"type:varchar(255)" validate:"required"`
	Content   string           `json:"content" gorm:"type:varchar(1024)" validate:"required"`
	Success   bool             `json:"success" gorm:"type:boolean"`
	Reason    string           `json:"reason" gorm:"type:varchar(255)"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `json:"deleted_at" gorm:"index"`
}
