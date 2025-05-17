package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

func (a authService) HandleRefreshAccessTokenRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	refreshToken := c.Locals("token").(*entity.JsonWebToken)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	token, err := a.jsonWebTokenRepository.GenerateAccessToken(user, "tongla.dev", "tongla.dev", "", "", "", refreshToken.ID)
	if err != nil {
		panic(err)
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
