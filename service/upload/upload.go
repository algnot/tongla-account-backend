package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"tongla-account/di/config"
)

type UploadService interface {
	HandleUploadFileRouter(c *fiber.Ctx) error
}

type uploadService struct {
	db     *gorm.DB
	config config.AppConfig
}

func ProvideUploadService(db *gorm.DB, config config.AppConfig) UploadService {
	return &uploadService{
		db:     db,
		config: config,
	}
}
