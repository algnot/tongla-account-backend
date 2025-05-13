package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
	"tongla-account/di/config"
	"tongla-account/entity"
	"tongla-account/repository"
)

func RequireAuth(db *gorm.DB, config config.AppConfig, tokenType entity.JsonTokenType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		encryptorRepository := repository.ProvideEncryptorRepository(db, config)
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		secretHash, err := encryptorRepository.GetPassphrase()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get secret hash",
			})
		}

		tokenParsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretHash.Hash, nil
		})

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid exp claim",
			})
		}

		if time.Now().Unix() > int64(exp) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}

		if err != nil || !tokenParsed.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		jsonWebTokenRepository := repository.ProvideJsonWebTokenRepository(db, config)
		jwtId := claims["sub"].(string)
		jwtEnt, err := jsonWebTokenRepository.GetTokenById(jwtId)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token Not Found",
			})
		}

		if jwtEnt.Type != tokenType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong token type",
			})
		}

		if jwtEnt.Revoked {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is revoked",
			})
		}

		accountRepository := repository.ProvideAccountRepository(db, config)
		user, err := accountRepository.FindById(jwtEnt.AccountId)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Account Not Found",
			})
		}

		c.Locals("user", user)
		c.Locals("token", jwtEnt)
		return c.Next()
	}
}
