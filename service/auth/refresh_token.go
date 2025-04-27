package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

func (a authService) HandleRefreshAccessTokenRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	token, err := a.jsonWebTokenRepository.GenerateAccessToken(user, "tongla.dev", "tongla.dev")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating access token",
		})
	}
	
	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
