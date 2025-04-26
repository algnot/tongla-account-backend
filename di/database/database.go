package database

import (
	"database/sql"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tongla-account/di/config"
)

func InitDatabase() (*gorm.DB, error) {
	DBConfig := getDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.DBName)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func getDatabaseConfig() config.DatabaseConfig {
	var app config.AppConfig
	envconfig.MustProcess("APP", &app.DatabaseConfig)
	return app.DatabaseConfig
}
