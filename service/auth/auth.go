package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/repository"
)

type AuthService interface {
	HandleRegisterRouter(c *fiber.Ctx) error
	HandleVerifyEmailRouter(c *fiber.Ctx) error
	HandleResendVerifyEmailRouter(c *fiber.Ctx) error
	HandleResendVerify2FARouter(c *fiber.Ctx) error
	HandleLoginRouter(c *fiber.Ctx) error
}

type authService struct {
	encryptorRepository repository.EncryptorRepository
	accountRepository   repository.AccountRepository
	tokenRepository     repository.TokenRepository
	db                  *gorm.DB
	config              config.AppConfig
}

func ProvideAuthService(db *gorm.DB, config config.AppConfig) AuthService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	accountRepository := repository.ProvideAccountRepository(db, config)
	tokenRepository := repository.ProvideTokenRepository(db, config)
	return &authService{
		db:                  db,
		config:              config,
		tokenRepository:     tokenRepository,
		accountRepository:   accountRepository,
		encryptorRepository: encryptorRepository,
	}
}
