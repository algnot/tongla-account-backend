package databasetest

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tongla-account/entity/migrater"
	"os"
)

func InitTestDatabase() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	con, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := con.Ping(); err != nil {
		panic(err)
	}

	err = migrater.AutoMigrate(db)
	if err != nil {
		return nil, nil
	}

	return db, func() {
		db.Rollback()
		err := con.Close()
		if err != nil {
			return
		}
		err = os.Remove("gorm.db")
		if err != nil {
			panic(err)
		}
	}
}
