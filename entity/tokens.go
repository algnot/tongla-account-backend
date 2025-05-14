package entity

import "time"

type TokenType string

const (
	TokenVerifyEmail TokenType = "verifyEmail"
	TokenLogin       TokenType = "login"
	TokenAuthCode    TokenType = "authCode"
)

type Token struct {
	ID        string    `gorm:"type:varchar(255);primarykey"`
	AccountID string    `gorm:"type:varchar(255);index"`
	Type      TokenType `gorm:"type:varchar(255);index"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex"`
	Ref       string    `gorm:"type:varchar(255)"`
	Used      bool      `gorm:"default:false"`
	ExpireAt  time.Time
	CreatedAt time.Time
}
