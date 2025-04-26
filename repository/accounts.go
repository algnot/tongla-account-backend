package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type AccountRepository interface {
	CreateAccount(account *entity.Account) (*entity.Account, error)
	UpdateAccount(account *entity.Account) (*entity.Account, error)
	FindByUsername(username string) (*entity.Account, error)
	FindByEmail(email string) (*entity.Account, error)
	FindById(id string) (*entity.Account, error)

	isDuplicateAccount(account *entity.Account) (bool, error)
	sendVerifyEmail(account *entity.Account) error
}

type accountRepository struct {
	db                     *gorm.DB
	config                 config.AppConfig
	encryptorRepository    EncryptorRepository
	notificationRepository NotificationRepository
	tokenRepository        TokenRepository
}

func (a accountRepository) UpdateAccount(account *entity.Account) (*entity.Account, error) {
	result := a.db.Updates(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (a accountRepository) FindById(id string) (*entity.Account, error) {
	var ent entity.Account
	result := a.db.First(&ent, "id = ?", id)
	if result.Error != nil {
		return &entity.Account{}, result.Error
	}
	return &ent, nil
}

func (a accountRepository) sendVerifyEmail(account *entity.Account) error {
	token, err := a.encryptorRepository.GeneratePassphrase(50)
	if err != nil {
		return err
	}

	tokenEnt := &entity.Token{
		AccountID: account.ID,
		Type:      entity.TokenVerifyEmail,
		Token:     token,
		ExpireAt:  time.Now().Add(30 * time.Minute),
	}

	tokenEnt, err = a.tokenRepository.CreateToken(tokenEnt)
	if err != nil {
		return err
	}

	err = a.notificationRepository.SendNotification(&entity.Notification{
		Type:  entity.NotificationEmail,
		Email: account.Email,
		Title: "Verify Tongla Account",
		Content: fmt.Sprintf(`Hello, %s

To verify your Tongla account, we need to confirm your email. Please click the following link or copy & paste into the browser:

%s/auth/verify-email?token=%s

The link is expird in 30 minutes.

Best regards,
Tongla
www.tongla.dev`, a.encryptorRepository.Decrypt(account.Username), a.config.ServerConfig.FrontendPath, tokenEnt.Token),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a accountRepository) FindByUsername(username string) (*entity.Account, error) {
	var ent entity.Account
	result := a.db.First(&ent, "username = ?", a.encryptorRepository.Encrypt(username))
	if result.Error != nil {
		return &entity.Account{}, result.Error
	}
	return &ent, nil
}

func (a accountRepository) FindByEmail(email string) (*entity.Account, error) {
	var ent entity.Account
	result := a.db.First(&ent, "email = ?", a.encryptorRepository.Encrypt(email))
	if result.Error != nil {
		return &entity.Account{}, result.Error
	}
	return &ent, nil
}

func (a accountRepository) isDuplicateAccount(account *entity.Account) (bool, error) {
	email := account.Email
	username := account.Username

	existingAccount, err := a.FindByUsername(a.encryptorRepository.Decrypt(username))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, err
		}
	}

	if existingAccount.ID != "" {
		return true, errors.New(fmt.Sprintf("account with username %s already exists",
			a.encryptorRepository.Decrypt(existingAccount.Username)))
	}

	existingAccount, err = a.FindByEmail(a.encryptorRepository.Decrypt(email))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, err
		}
	}

	if existingAccount.ID != "" {
		return true, errors.New(fmt.Sprintf("account with email %s already exists",
			a.encryptorRepository.Decrypt(existingAccount.Email)))
	}

	return false, nil
}

func (a accountRepository) CreateAccount(account *entity.Account) (*entity.Account, error) {
	_, err := a.isDuplicateAccount(account)
	if err != nil {
		return nil, err
	}

	id, err := a.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	account.ID = id

	result := a.db.Create(account)

	if result.Error != nil {
		return nil, result.Error
	}

	err = a.sendVerifyEmail(account)
	if err != nil {
		return account, err
	}

	return account, nil
}

func ProvideAccountRepository(db *gorm.DB, config config.AppConfig) AccountRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	notificationRepository := ProvideNotificationRepository(db, config)
	tokenRepository := ProvideTokenRepository(db, config)
	return &accountRepository{
		db:                     db,
		config:                 config,
		encryptorRepository:    encryptorRepository,
		notificationRepository: notificationRepository,
		tokenRepository:        tokenRepository,
	}
}
