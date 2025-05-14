package entity

import (
	"gorm.io/gorm"
	"time"
)

type Service struct {
	Name         string         `json:"name" gorm:"type:varchar(255)"`
	ClientId     string         `json:"client_id" gorm:"type:varchar(255);primarykey"`
	ClientSecret string         `json:"client_secret" gorm:"type:varchar(255)"`
	RedirectUri  string         `json:"redirect_uri" gorm:"type:varchar(255)"`
	Issuer       string         `json:"issuer" gorm:"type:varchar(255)"`
	Owner        string         `json:"owner" gorm:"type:varchar(255);index"`
	Scope        string         `json:"scope" gorm:"type:varchar(255)"`
	GrantType    string         `json:"grant_type" gorm:"type:varchar(255)"`
	ResponseType string         `json:"response_type" gorm:"type:varchar(255)"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
