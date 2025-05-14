package repository

import (
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type ServiceRepository interface {
	GetByClientId(clientId string) (*entity.Service, error)
}

type serviceRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (s serviceRepository) GetByClientId(clientId string) (*entity.Service, error) {
	var ent entity.Service
	result := s.db.First(&ent, "client_id = ?", clientId)
	if result.Error != nil {
		return &entity.Service{}, result.Error
	}
	return &ent, nil
}

func ProvideServiceRepository(db *gorm.DB, config config.AppConfig) ServiceRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &serviceRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
