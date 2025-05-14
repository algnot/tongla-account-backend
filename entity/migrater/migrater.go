package migrater

import (
	"gorm.io/gorm"
	"tongla-account/entity"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.ApiKeys{},
		&entity.Encryptor{},
		&entity.Account{},
		&entity.Notification{},
		&entity.Token{},
		&entity.JsonWebToken{},
		&entity.Service{})
	if err != nil {
		return err
	}
	return nil
}
