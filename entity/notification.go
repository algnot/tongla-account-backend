package entity

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type NotificationType string

const (
	NotificationEmail NotificationType = "email"
	NotificationWeb   NotificationType = "web"
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

func (n *Notification) ToResponse() fiber.Map {
	return fiber.Map{
		"title":   n.Title,
		"content": n.Content,
		"reason":  n.Reason,
		"success": n.Success,
		"created": n.CreatedAt,
	}
}
