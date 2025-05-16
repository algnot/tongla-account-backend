package router

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"tongla-account/di/config"
	"tongla-account/di/database"
	"tongla-account/entity"
	"tongla-account/service/api_keys"
	service2 "tongla-account/service/auth"
	"tongla-account/service/middleware"
	service4 "tongla-account/service/notification"
	service3 "tongla-account/service/open_id"
)

func InitRouter(server *fiber.App) {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	appConfig := config.GetConfig()
	apiKeysService := service.ProvideApiKeysService(db, appConfig)
	authService := service2.ProvideAuthService(db, appConfig)
	openIdServer := service3.ProvideOpenIdServiceService(db, appConfig)
	notificationServer := service4.ProvideNotificationServiceService(db, appConfig)

	server.Post("/secret/generate", apiKeysService.HandleSecretPostRouter)
	server.Post("/secret/rotate", apiKeysService.HandleRotatePostRouter)
	server.Get("/secret/verify", apiKeysService.HandleVerifyGetRouter)

	server.Post("/auth/register", authService.HandleRegisterRouter)
	server.Post("/auth/verify-email", authService.HandleVerifyEmailRouter)
	server.Post("/auth/verify-2FA", authService.HandleResendVerify2FARouter)
	server.Post("/auth/resend-verify-email", authService.HandleVerifyEmailRouter)
	server.Post("/auth/login", authService.HandleLoginRouter)
	server.Post("/auth/login-with-code", authService.HandleLoginWithCodeRouter)
	server.Post("/auth/request-email-login", authService.HandleRequestLoginWithEmailRouter)
	server.Post("/auth/login-with-token", authService.HandleLoginWithTokenRouter)

	refreshProtected := server.Group("/auth/refresh", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenRefreshToken))
	refreshProtected.Post("/", authService.HandleRefreshAccessTokenRouter)

	refreshProtected2 := server.Group("/auth/logout", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenRefreshToken))
	refreshProtected2.Post("/", authService.HandleLogoutRouter)

	authProtected := server.Group("/account", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenAccessToken))
	authProtected.Get("/me", authService.HandleGetUserInfoRouter)
	authProtected.Put("/update-user", authService.HandleUpdateUserRouter)
	authProtected.Get("/all-device", authService.HandleGetAllDeviceRouter)
	authProtected.Delete("/delete-device", authService.HandleDeleteDeviceRouter)
	authProtected.Get("/all-service", authService.HandleGetAllServiceRouter)

	openIdProtected := server.Group("/openid/get-service", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenRefreshToken))
	openIdProtected.Post("/", openIdServer.HandleGetServiceRouter)

	openIdGetUserInfoProtected := server.Group("/openid/userinfo", middleware.RequireOpenIdAuth(db, appConfig))
	openIdGetUserInfoProtected.Get("/", openIdServer.HandleGetUserInfoRouter)

	server.Post("/openid/token", openIdServer.HandleGetTokenRouter)
	server.Get("/openid/.well-known/configuration", openIdServer.HandleCertificateRouter)
	server.Get("/openid/.well-known/jwks.json", openIdServer.HandleJWKSRouter)

	notificationProtected := server.Group("/notification", middleware.RequireAuth(db, appConfig, entity.JsonWebTokenAccessToken))
	notificationProtected.Get("/get-all", notificationServer.HandleGetAllNotificationsRouter)
}
