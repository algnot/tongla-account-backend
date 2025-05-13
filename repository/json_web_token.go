package repository

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type JsonWebTokenRepository interface {
	GenerateToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string) (*entity.JwtTokenResponse, error)
	GetTokenById(jwtId string) (*entity.JsonWebToken, error)
	GenerateAccessToken(user *entity.Account, issuer string, audience string, ref string) (string, error)

	createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType, user *entity.Account, ref string) (*entity.JsonWebToken, error)
	generateRefreshToken(user *entity.Account, issuer string, audience string, userAgent string, deviceID string) (string, *entity.JsonWebToken, error)
}

type jsonWebTokenRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (j jsonWebTokenRepository) GetTokenById(jwtId string) (*entity.JsonWebToken, error) {
	var ent entity.JsonWebToken
	result := j.db.First(&ent, "id = ?", jwtId)
	if result.Error != nil {
		return &entity.JsonWebToken{}, result.Error
	}
	return &ent, nil
}

func (j jsonWebTokenRepository) createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType, user *entity.Account, ref string) (*entity.JsonWebToken, error) {
	ent := &entity.JsonWebToken{
		ID:        token.Sub,
		AccountId: user.ID,
		Type:      tokenType,
		Iat:       token.Iat,
		Exp:       token.Exp,
		Issuer:    token.Iss,
		Audience:  token.Aud,
		UserAgent: token.UserAgent,
		DeviceID:  token.DeviceID,
		Ref:       ref,
	}

	result := j.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	return ent, nil
}

func (j jsonWebTokenRepository) GenerateToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string) (*entity.JwtTokenResponse, error) {
	refreshToken, refreshTokenEnt, err := j.generateRefreshToken(userEnt, issuer, audience, userAgent, deviceID)
	if err != nil {
		return nil, err
	}

	accessToken, err := j.GenerateAccessToken(userEnt, issuer, audience, refreshTokenEnt.ID)
	if err != nil {
		return nil, err
	}

	return &entity.JwtTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j jsonWebTokenRepository) generateRefreshToken(user *entity.Account, issuer string, audience string, userAgent string, deviceID string) (string, *entity.JsonWebToken, error) {
	passphrase, err := j.encryptorRepository.GetPassphrase()
	if err != nil {
		return "", nil, err
	}

	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", nil, err
	}

	secretKey := passphrase.Hash
	jwtEnt := &entity.JwtToken{
		Sub:       id,
		Iat:       time.Now().Unix(),
		Exp:       time.Now().Add(time.Minute * 60 * 24 * 7).Unix(),
		Iss:       issuer,
		Aud:       audience,
		UserAgent: userAgent,
		DeviceID:  deviceID,
	}

	result, err := j.createJsonWebToken(jwtEnt, entity.JsonWebTokenRefreshToken, user, "")
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtEnt.ToMapClaims())
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil, err
	}

	return signedToken, result, nil
}

func (j jsonWebTokenRepository) GenerateAccessToken(user *entity.Account, issuer string, audience string, ref string) (string, error) {
	passphrase, err := j.encryptorRepository.GetPassphrase()
	if err != nil {
		return "", err
	}

	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", err
	}

	secretKey := passphrase.Hash
	jwtEnt := &entity.JwtToken{
		Sub:       id,
		Iat:       time.Now().Unix(),
		Exp:       time.Now().Add(time.Minute * 15).Unix(),
		Iss:       issuer,
		Aud:       audience,
		UserAgent: "",
		DeviceID:  "",
	}

	_, err = j.createJsonWebToken(jwtEnt, entity.JsonWebTokenAccessToken, user, ref)
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
