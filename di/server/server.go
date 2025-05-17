package server

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
	"tongla-account/di/config"
	router "tongla-account/service"
)

func InitApiServer() error {
	appConfig := config.GetConfig()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              appConfig.ServerConfig.SentryDns,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	sentryHandler := sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	})

	serverConfig := getServerConfig()
	server := fiber.New()
	server.Use(sentryHandler)
	server.Use(func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				if hub := sentryfiber.GetHubFromContext(c); hub != nil {
					hub.Recover(r)
					hub.Flush(1 * time.Second)
				}

				errMsg := fmt.Sprintf("%v", r)
				err = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": errMsg,
				})
			}
		}()
		return c.Next()
	})
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
