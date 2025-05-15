package repository

import (
	"errors"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"tongla-account/di/config"
	"tongla-account/entity"
)

type NotificationRepository interface {
	SendNotification(notificationEnt *entity.Notification) error
	CreateNotification(notificationEnt *entity.Notification) (*entity.Notification, error)
	UpdateNotification(notificationEnt *entity.Notification) (*entity.Notification, error)

	sendEmail(notificationEnt *entity.Notification) error
}

type notificationRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (n notificationRepository) UpdateNotification(notificationEnt *entity.Notification) (*entity.Notification, error) {
	result := n.db.Model(notificationEnt).Updates(notificationEnt)
	if result.Error != nil {
		return nil, result.Error
	}
	return notificationEnt, nil
}

func (n notificationRepository) CreateNotification(notificationEnt *entity.Notification) (*entity.Notification, error) {
	id, err := n.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	notificationEnt.ID = id

	result := n.db.Create(notificationEnt)
	if result.Error != nil {
		return nil, result.Error
	}
	return notificationEnt, nil
}

func (n notificationRepository) sendEmail(notificationEnt *entity.Notification) error {
	m := gomail.NewMessage()
	m.SetHeader("From", n.config.EmailConfig.Sender)
	m.SetHeader("To", n.encryptorRepository.Decrypt(notificationEnt.Email))
	m.SetHeader("Subject", notificationEnt.Title)
	m.SetBody("text/plain", notificationEnt.Content)

	d := gomail.NewDialer(n.config.EmailConfig.Host,
		n.config.EmailConfig.Port,
		n.config.EmailConfig.Sender,
		n.config.EmailConfig.Password)

	if err := d.DialAndSend(m); err != nil {
		notificationEnt.Success = false
		notificationEnt.Reason = err.Error()
		_, err := n.UpdateNotification(notificationEnt)
		if err != nil {
			return err
		}
		return err
	}

	notificationEnt.Success = true
	notificationEnt.Reason = ""
	_, err := n.UpdateNotification(notificationEnt)
	if err != nil {
		return err
	}

	return nil
}

func (n notificationRepository) SendNotification(notificationEnt *entity.Notification) error {
	createdNotification, err := n.CreateNotification(notificationEnt)
	if err != nil {
		return err
	}

	if createdNotification.Type == entity.NotificationEmail {
		return n.sendEmail(createdNotification)
	}

	if createdNotification.Type == entity.NotificationWeb {
		notificationEnt.Success = true
		_, err := n.UpdateNotification(notificationEnt)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("invalid notification type")
}

func ProvideNotificationRepository(db *gorm.DB, config config.AppConfig) NotificationRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &notificationRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
