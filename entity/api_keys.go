package entity

import (
	"gorm.io/gorm"
	"time"
)

type ApiKeysRequest struct {
	Name string `json:"name" validate:"required"`
}

type ApiKeys struct {
	ID        string         `json:"id" gorm:"type:varchar(255);primarykey"`
	Name      string         `json:"name" gorm:"type:varchar(255)" validate:"required"`
	Secret    EncryptedField `json:"secret" gorm:"type:varbinary(512)" validate:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
