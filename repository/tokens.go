package repository

import (
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type TokenRepository interface {
	CreateToken(*entity.Token) (*entity.Token, error)
	UpdateToken(*entity.Token) (*entity.Token, error)
	FindKeyByToken(token string) (*entity.Token, error)
}

type tokenRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (t tokenRepository) CreateToken(token *entity.Token) (*entity.Token, error) {
	id, err := t.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	token.ID = id

	result := t.db.Create(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (t tokenRepository) FindKeyByToken(token string) (*entity.Token, error) {
	var ent *entity.Token
	result := t.db.First(&ent, "token = ?", token)
	if result.Error != nil {
		return &entity.Token{}, result.Error
	}
	return ent, nil
}

func (t tokenRepository) UpdateToken(token *entity.Token) (*entity.Token, error) {
	result := t.db.Updates(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func ProvideTokenRepository(db *gorm.DB, config config.AppConfig) TokenRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &tokenRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
