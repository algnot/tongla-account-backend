package entity

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID         string         `json:"id" gorm:"type:varchar(255);primarykey"`
	Username   EncryptedField `json:"username" gorm:"type:varbinary(512)" validate:"required"`
	Email      EncryptedField `json:"email" gorm:"type:varbinary(512)" validate:"required"`
	Firstname  EncryptedField `json:"firstname" gorm:"type:varbinary(512)" validate:"required"`
	Lastname   EncryptedField `json:"lastname" gorm:"type:varbinary(512)" validate:"required"`
	IsVerified bool           `json:"isVerified" gorm:"type:boolean;default:false"`
	Secret     EncryptedField `json:"secret" gorm:"type:varbinary(512)" validate:"required"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type AccountResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	IsVerified bool   `json:"isVerified"`
}

func (a *Account) ToResponse(decrypt func(EncryptedField) string) AccountResponse {
	return AccountResponse{
		ID:         a.ID,
		Username:   decrypt(a.Username),
		Email:      decrypt(a.Email),
		Firstname:  decrypt(a.Firstname),
		Lastname:   decrypt(a.Lastname),
		IsVerified: a.IsVerified,
	}
}
