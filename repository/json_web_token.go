package repository

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type JsonWebTokenRepository interface {
	GenerateToken(userEnt *entity.Account, issuer string, audience string) (*entity.JwtTokenResponse, error)

	createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType) (*entity.JsonWebToken, error)
	generateRefreshToken(user *entity.Account, issuer string, audience string) (string, error)
	generateAccessToken(user *entity.Account, issuer string, audience string) (string, error)
}

type jsonWebTokenRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (j jsonWebTokenRepository) createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType) (*entity.JsonWebToken, error) {
	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}

	ent := &entity.JsonWebToken{
		ID:        id,
		AccountId: token.Sub,
		Type:      tokenType,
		Iat:       token.Iat,
		Exp:       token.Exp,
		Issuer:    token.Iss,
		Audience:  token.Aud,
	}

	result := j.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	return ent, nil
}

func (j jsonWebTokenRepository) GenerateToken(userEnt *entity.Account, issuer string, audience string) (*entity.JwtTokenResponse, error) {
	refreshToken, err := j.generateRefreshToken(userEnt, issuer, audience)
	if err != nil {
		return nil, err
	}

	accessToken, err := j.generateAccessToken(userEnt, issuer, audience)
	if err != nil {
		return nil, err
	}

	return &entity.JwtTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j jsonWebTokenRepository) generateRefreshToken(user *entity.Account, issuer string, audience string) (string, error) {
	passphrase, err := j.encryptorRepository.GetPassphrase()
	if err != nil {
		return "", err
	}

	secretKey := passphrase.Hash
	jwtEnt := &entity.JwtToken{
		Sub: user.ID,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Minute * 60 * 24 * 7).Unix(),
		Iss: issuer,
		Aud: audience,
	}

	_, err = j.createJsonWebToken(jwtEnt, entity.JsonWebTokenRefreshToken)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtEnt.ToMapClaims())
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j jsonWebTokenRepository) generateAccessToken(user *entity.Account, issuer string, audience string) (string, error) {
	passphrase, err := j.encryptorRepository.GetPassphrase()
	if err != nil {
		return "", err
	}

	secretKey := passphrase.Hash
	jwtEnt := &entity.JwtToken{
		Sub: user.ID,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Minute * 15).Unix(),
		Iss: issuer,
		Aud: audience,
	}

	_, err = j.createJsonWebToken(jwtEnt, entity.JsonWebTokenAccessToken)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtEnt.ToMapClaims())
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ProvideJsonWebTokenRepository(db *gorm.DB, config config.AppConfig) JsonWebTokenRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &jsonWebTokenRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
