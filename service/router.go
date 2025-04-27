package router

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"tongla-account/di/config"
	"tongla-account/di/database"
	"tongla-account/service/api_keys"
	service2 "tongla-account/service/auth"
)

func InitRouter(server *fiber.App) {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	appConfig := config.GetConfig()
	apiKeysService := service.ProvideApiKeysService(db, appConfig)
	authService := service2.ProvideAuthService(db, appConfig)

	server.Post("/secret/generate", apiKeysService.HandleSecretPostRouter)
	server.Post("/secret/rotate", apiKeysService.HandleRotatePostRouter)
	server.Get("/secret/verify", apiKeysService.HandleVerifyGetRouter)

	server.Post("/auth/register", authService.HandleRegisterRouter)
	server.Post("/auth/verify-email", authService.HandleVerifyEmailRouter)
	server.Post("/auth/verify-2FA", authService.HandleResendVerify2FARouter)
	server.Post("/auth/resend-verify-email", authService.HandleVerifyEmailRouter)
	server.Post("/auth/login", authService.HandleLoginRouter)
	server.Post("/auth/login-with-code", authService.HandleLoginWithCodeRouter)
}
