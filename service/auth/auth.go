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
	HandleLoginWithCodeRouter(c *fiber.Ctx) error
	HandleLoginWithTokenRouter(c *fiber.Ctx) error
	HandleRequestLoginWithEmailRouter(c *fiber.Ctx) error
	HandleGetUserInfoRouter(c *fiber.Ctx) error
	HandleRefreshAccessTokenRouter(c *fiber.Ctx) error
	HandleUpdateUserRouter(c *fiber.Ctx) error
	HandleGetAllDeviceRouter(c *fiber.Ctx) error
	HandleLogoutRouter(c *fiber.Ctx) error
	HandleDeleteDeviceRouter(c *fiber.Ctx) error
	HandleGetAllServiceRouter(c *fiber.Ctx) error
}

type authService struct {
	encryptorRepository    repository.EncryptorRepository
	accountRepository      repository.AccountRepository
	tokenRepository        repository.TokenRepository
	jsonWebTokenRepository repository.JsonWebTokenRepository
	notificationRepository repository.NotificationRepository
	serviceRepository      repository.ServiceRepository
	db                     *gorm.DB
	config                 config.AppConfig
}

func ProvideAuthService(db *gorm.DB, config config.AppConfig) AuthService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	accountRepository := repository.ProvideAccountRepository(db, config)
	tokenRepository := repository.ProvideTokenRepository(db, config)
	jsonWebTokenRepository := repository.ProvideJsonWebTokenRepository(db, config)
	notificationRepository := repository.ProvideNotificationRepository(db, config)
	serviceRepository := repository.ProvideServiceRepository(db, config)
	return &authService{
		db:                     db,
		config:                 config,
		tokenRepository:        tokenRepository,
		accountRepository:      accountRepository,
		encryptorRepository:    encryptorRepository,
		jsonWebTokenRepository: jsonWebTokenRepository,
		notificationRepository: notificationRepository,
		serviceRepository:      serviceRepository,
	}
}
