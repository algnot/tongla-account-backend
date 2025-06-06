package repository

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
	"tongla-account/di/config"
	"tongla-account/entity"
	"tongla-account/util"
)

type JsonWebTokenRepository interface {
	GenerateToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string, clientId string) (*entity.JwtTokenResponse, error)
	GetTokenById(jwtId string) (*entity.JsonWebToken, error)
	GenerateAccessToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string, clientId string, ref string) (string, error)
	GetAllActiveTokenByAccountId(userId string, tokenType entity.JsonTokenType) (*[]entity.JsonWebToken, error)
	RevokedAllActiveTokenByRefId(refId string) error
	GetActiveRefreshTokenByClientId(clientId string) (*entity.JsonWebToken, error)

	createJsonWebToken(token *entity.JwtToken, tokenType entity.JsonTokenType, user *entity.Account, ref string) (*entity.JsonWebToken, error)
	generateRefreshToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string, clientId string) (string, *entity.JsonWebToken, error)
}

type jsonWebTokenRepository struct {
	db                     *gorm.DB
	config                 config.AppConfig
	encryptorRepository    EncryptorRepository
	notificationRepository NotificationRepository
}

func (j jsonWebTokenRepository) GetActiveRefreshTokenByClientId(clientId string) (*entity.JsonWebToken, error) {
	var token entity.JsonWebToken
	now := time.Now().Unix()

	result := j.db.
		Where("client_id = ? AND type = ? AND exp > ? AND revoked = 0", clientId, entity.JsonWebTokenRefreshToken, now).
		Order("exp DESC").
		First(&token)

	if result.Error != nil {
		return nil, result.Error
	}

	return &token, nil
}

func (j jsonWebTokenRepository) RevokedAllActiveTokenByRefId(refId string) error {
	result := j.db.
		Model(&entity.JsonWebToken{}).
		Where("revoked = 0 AND (id = ? OR ref = ?)", refId, refId).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (j jsonWebTokenRepository) GetAllActiveTokenByAccountId(userId string, tokenType entity.JsonTokenType) (*[]entity.JsonWebToken, error) {
	var tokens []entity.JsonWebToken
	now := time.Now().Unix()

	result := j.db.
		Where("account_id = ? AND exp > ? AND type = ? AND revoked = 0", userId, now, tokenType).
		Find(&tokens)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tokens, nil
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
		ClientId:  token.ClientId,
		Ref:       ref,
	}

	result := j.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	return ent, nil
}

func (j jsonWebTokenRepository) GenerateToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string, clientId string) (*entity.JwtTokenResponse, error) {
	refreshToken, refreshTokenEnt, err := j.generateRefreshToken(userEnt, issuer, audience, userAgent, deviceID, clientId)
	if err != nil {
		return nil, err
	}

	accessToken, err := j.GenerateAccessToken(userEnt, issuer, audience, userAgent, deviceID, clientId, refreshTokenEnt.ID)
	if err != nil {
		return nil, err
	}

	_ = j.notificationRepository.SendNotification(&entity.Notification{
		Type:    entity.NotificationWeb,
		Email:   userEnt.Email,
		Title:   "New login to tongla account",
		Content: fmt.Sprintf(util.GetWebNotificationContent("login"), issuer, deviceID),
		Reason:  "login",
	})

	return &entity.JwtTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j jsonWebTokenRepository) generateRefreshToken(user *entity.Account, issuer string, audience string, userAgent string, deviceID string, clientId string) (string, *entity.JsonWebToken, error) {
	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", nil, err
	}

	jwtEnt := &entity.JwtToken{
		Sub:       id,
		Iat:       time.Now().Unix(),
		Exp:       time.Now().Add(time.Minute * 60 * 24 * 7).Unix(),
		Iss:       issuer,
		Aud:       audience,
		UserAgent: userAgent,
		DeviceID:  deviceID,
		ClientId:  clientId,
		Email:     j.encryptorRepository.Decrypt(user.Email),
		Profile:   j.encryptorRepository.Decrypt(user.Firstname) + " " + j.encryptorRepository.Decrypt(user.Lastname),
	}

	result, err := j.createJsonWebToken(jwtEnt, entity.JsonWebTokenRefreshToken, user, "")
	if err != nil {
		return "", nil, err
	}

	privateKey, err := util.EnsureRSAKeyPair()
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtEnt.ToMapClaims())
	token.Header["kid"] = "tongla.dev"
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", nil, err
	}

	return signedToken, result, nil
}

func (j jsonWebTokenRepository) GenerateAccessToken(userEnt *entity.Account, issuer string, audience string, userAgent string, deviceID string, clientId string, ref string) (string, error) {
	id, err := j.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return "", err
	}

	jwtEnt := &entity.JwtToken{
		Sub:       id,
		Iat:       time.Now().Unix(),
		Exp:       time.Now().Add(time.Minute * 15).Unix(),
		Iss:       issuer,
		Aud:       audience,
		Email:     j.encryptorRepository.Decrypt(userEnt.Email),
		Profile:   j.encryptorRepository.Decrypt(userEnt.Firstname) + " " + j.encryptorRepository.Decrypt(userEnt.Lastname),
		UserAgent: userAgent,
		DeviceID:  deviceID,
		ClientId:  clientId,
	}

	_, err = j.createJsonWebToken(jwtEnt, entity.JsonWebTokenAccessToken, userEnt, ref)
	if err != nil {
		return "", err
	}

	privateKey, err := util.EnsureRSAKeyPair()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtEnt.ToMapClaims())
	token.Header["kid"] = "tongla.dev"
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ProvideJsonWebTokenRepository(db *gorm.DB, config config.AppConfig) JsonWebTokenRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	notificationRepository := ProvideNotificationRepository(db, config)
	return &jsonWebTokenRepository{
		db:                     db,
		config:                 config,
		encryptorRepository:    encryptorRepository,
		notificationRepository: notificationRepository,
	}
}
