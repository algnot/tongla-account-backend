package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

func (a authService) HandleLogoutRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	refreshToken := c.Locals("token").(*entity.JsonWebToken)
	if refreshToken == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Refresh token missing",
		})
	}

	err := a.jsonWebTokenRepository.RevokedAllActiveTokenByRefId(refreshToken.ID)
	if err != nil {
		panic(err)
	}

	return c.JSON(fiber.Map{})
}
