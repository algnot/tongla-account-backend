package ditest

import (
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/repository"
	"tongla-account/service/api_keys"
)

type ApplicationTestSuite struct {
	DB                  *gorm.DB
	Config              config.AppConfig
	ApiKeysRepository   repository.ApiKeysRepository
	EncryptorRepository repository.EncryptorRepository
	ApiKeysService      service.ApiKeysService
}
