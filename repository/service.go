package repository

import (
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type ServiceRepository interface {
	GetByClientId(clientId string) (*entity.Service, error)
	GetAllServiceByAccountId(accountId string) ([]*entity.Service, error)
	CreateService(service *entity.Service) (*entity.Service, error)
}

type serviceRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (s serviceRepository) CreateService(service *entity.Service) (*entity.Service, error) {
	clientId, err := s.encryptorRepository.GeneratePassphrase(30)
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.encryptorRepository.GeneratePassphrase(50)
	if err != nil {
		return nil, err
	}

	service.ClientId = clientId
	service.ClientSecret = clientSecret

	result := s.db.Create(service)

	if result.Error != nil {
		return nil, result.Error
	}

	return service, nil
}

func (s serviceRepository) GetAllServiceByAccountId(accountId string) ([]*entity.Service, error) {
	var services []*entity.Service

	result := s.db.
		Where("owner = ?", accountId).
		Find(&services)

	if result.Error != nil {
		return nil, result.Error
	}

	return services, nil
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
