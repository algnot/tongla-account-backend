package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/repository"
)

type OpenIdService interface {
	HandleGetServiceRouter(c *fiber.Ctx) error
	HandleGetTokenRouter(c *fiber.Ctx) error
	HandleCertificateRouter(c *fiber.Ctx) error
	HandleJWKSRouter(c *fiber.Ctx) error
	HandleGetUserInfoRouter(c *fiber.Ctx) error
	HandleAddServiceRouter(c *fiber.Ctx) error
}

type openIdService struct {
	encryptorRepository    repository.EncryptorRepository
	accountRepository      repository.AccountRepository
	tokenRepository        repository.TokenRepository
	jsonWebTokenRepository repository.JsonWebTokenRepository
	serviceRepository      repository.ServiceRepository
	db                     *gorm.DB
	config                 config.AppConfig
}

func ProvideOpenIdServiceService(db *gorm.DB, config config.AppConfig) OpenIdService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	accountRepository := repository.ProvideAccountRepository(db, config)
	tokenRepository := repository.ProvideTokenRepository(db, config)
	jsonWebTokenRepository := repository.ProvideJsonWebTokenRepository(db, config)
	serviceRepository := repository.ProvideServiceRepository(db, config)
	return &openIdService{
		db:                     db,
		config:                 config,
		tokenRepository:        tokenRepository,
		accountRepository:      accountRepository,
		encryptorRepository:    encryptorRepository,
		jsonWebTokenRepository: jsonWebTokenRepository,
		serviceRepository:      serviceRepository,
	}
}
