package entity

import (
	"gorm.io/gorm"
	"time"
)

type GenderType string

const (
	GenderMale     GenderType = "male"
	GenderFemale   GenderType = "female"
	GenderNotToSay GenderType = "notToSay"
)

type Account struct {
	ID         string         `json:"id" gorm:"type:varchar(255);primarykey"`
	Username   EncryptedField `json:"username" gorm:"type:varbinary(512)" validate:"required"`
	Email      EncryptedField `json:"email" gorm:"type:varbinary(512)" validate:"required"`
	Firstname  EncryptedField `json:"firstname" gorm:"type:varbinary(512)" validate:"required"`
	Lastname   EncryptedField `json:"lastname" gorm:"type:varbinary(512)" validate:"required"`
	IsVerified bool           `json:"isVerified" gorm:"type:boolean;default:false"`
	Secret     EncryptedField `json:"secret" gorm:"type:varbinary(512)" validate:"required"`
	Phone      EncryptedField `json:"phone" gorm:"type:varbinary(512)"`
	Address    EncryptedField `json:"address" gorm:"type:varbinary(512)"`
	Birthdate  *time.Time     `json:"birthdate" gorm:"type:date"`
	Gender     GenderType     `json:"gender" gorm:"type:varchar(20);default:'notToSay'"`
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
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Birthdate  string `json:"birthdate"`
	Gender     string `json:"gender"`
	IsVerified bool   `json:"is_verified"`
}

func (a *Account) ToResponse(decrypt func(EncryptedField) string) AccountResponse {
	var birthdateStr string
	if a.Birthdate == nil {
		birthdateStr = ""
	} else if a.Birthdate.IsZero() {
		birthdateStr = ""
	} else {
		birthdateStr = a.Birthdate.Format(time.DateOnly)
	}

	return AccountResponse{
		ID:         a.ID,
		Username:   decrypt(a.Username),
		Email:      decrypt(a.Email),
		Firstname:  decrypt(a.Firstname),
		Lastname:   decrypt(a.Lastname),
		Phone:      decrypt(a.Phone),
		Address:    decrypt(a.Address),
		Birthdate:  birthdateStr,
		Gender:     string(a.Gender),
		IsVerified: a.IsVerified,
	}
}
