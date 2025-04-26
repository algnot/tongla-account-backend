package di

import (
	"tongla-account/di/server"
)

func InitApplication() error {
	err := server.InitApiServer()
	if err != nil {
		return err
	}
	return nil
}
