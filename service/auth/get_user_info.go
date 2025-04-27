package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

func (a authService) HandleGetUserInfoRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user.ToResponse(a.encryptorRepository.Decrypt))
}
