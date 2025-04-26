package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kelseyhightower/envconfig"
	"tongla-account/di/config"
	"tongla-account/service"
	"log"
)

func InitApiServer() error {
	serverConfig := getServerConfig()
	server := fiber.New()
	server.Use(recover2.New())
	server.Use(cors.New())

	router.InitRouter(server)

	log.Fatal(server.Listen(":" + serverConfig.Port))
	return nil
}

func getServerConfig() config.ServerConfig {
	var app config.AppConfig
	envconfig.MustProcess("APP", &app.ServerConfig)
	return app.ServerConfig
}
