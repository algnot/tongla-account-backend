package repository

import (
	"errors"
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type ApiKeysRepository interface {
	CreateKeyByName(name string) (string, error)
	FindKeyByName(name string) (entity.ApiKeys, error)
	RotateKeyByName(name string) (string, error)
	VerifyKey(key string) (entity.ApiKeys, error)
}

type apiKeysRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (a apiKeysRepository) CreateKeyByName(name string) (string, error) {
	existingKey, err := a.FindKeyByName(name)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	if existingKey.ID != "" {
		return "", errors.New("key already exists")
	}

	newKey, err := a.encryptorRepository.GeneratePassphrase(32)
	if err != nil {
		return "", err
	}

	id, err := a.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", err
	}

	encryptedKey := a.encryptorRepository.Encrypt(newKey)
	result := a.db.Create(&entity.ApiKeys{
		ID:     id,
		Name:   name,
		Secret: encryptedKey,
	})

	if result.Error != nil {
		return "", result.Error
	}

	return newKey, nil
}

func (a apiKeysRepository) FindKeyByName(name string) (entity.ApiKeys, error) {
	var apiKey entity.ApiKeys
	result := a.db.First(&apiKey, "name = ?", name)
	if result.Error != nil {
		return entity.ApiKeys{}, result.Error
	}
	return apiKey, nil
}

func (a apiKeysRepository) VerifyKey(key string) (entity.ApiKeys, error) {
	var apiKey entity.ApiKeys
	encryptedKey := a.encryptorRepository.Encrypt(key)
	result := a.db.First(&apiKey, "secret = ?", encryptedKey)
	if result.Error != nil {
		return entity.ApiKeys{}, result.Error
	}
	return apiKey, nil
}

func (a apiKeysRepository) RotateKeyByName(name string) (string, error) {
	passphrase, err := a.encryptorRepository.GeneratePassphrase(32)
	if err != nil {
		return "", err
	}

	existingKey, err := a.FindKeyByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("key not found please generate new one")
		}
		return "", err
	}

	existingKey.Secret = a.encryptorRepository.Encrypt(passphrase)

	result := a.db.Model(&entity.ApiKeys{}).Where("name = ?", name).Update("secret", a.encryptorRepository.Encrypt(passphrase))
	if result.Error != nil {
		return "", result.Error
	}

	return passphrase, nil
}

func ProvideApiKeysRepository(db *gorm.DB, config config.AppConfig) ApiKeysRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &apiKeysRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
