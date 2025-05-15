package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/repository"
)

type NotificationService interface {
	HandleGetAllNotificationsRouter(c *fiber.Ctx) error
}

type notificationService struct {
	encryptorRepository    repository.EncryptorRepository
	accountRepository      repository.AccountRepository
	tokenRepository        repository.TokenRepository
	jsonWebTokenRepository repository.JsonWebTokenRepository
	serviceRepository      repository.ServiceRepository
	notificationRepository repository.NotificationRepository
	db                     *gorm.DB
	config                 config.AppConfig
}

func ProvideNotificationServiceService(db *gorm.DB, config config.AppConfig) NotificationService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	accountRepository := repository.ProvideAccountRepository(db, config)
	tokenRepository := repository.ProvideTokenRepository(db, config)
	jsonWebTokenRepository := repository.ProvideJsonWebTokenRepository(db, config)
	serviceRepository := repository.ProvideServiceRepository(db, config)
	notificationRepository := repository.ProvideNotificationRepository(db, config)
	return &notificationService{
		db:                     db,
		config:                 config,
		tokenRepository:        tokenRepository,
		accountRepository:      accountRepository,
		encryptorRepository:    encryptorRepository,
		jsonWebTokenRepository: jsonWebTokenRepository,
		serviceRepository:      serviceRepository,
		notificationRepository: notificationRepository,
	}
}
