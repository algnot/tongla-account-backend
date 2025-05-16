package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

func (o openIdService) HandleGetUserInfoRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	token := c.Locals("token").(*entity.JsonWebToken)
	if token == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Refresh token missing",
		})
	}

	client, err := o.serviceRepository.GetByClientId(token.ClientId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Not found client",
		})
	}

	if client == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Not found client",
		})
	}

	user.Phone = o.encryptorRepository.Encrypt("")
	user.Address = o.encryptorRepository.Encrypt("")
	user.Gender = ""

	return c.Status(fiber.StatusOK).JSON(user.ToResponse(o.encryptorRepository.Decrypt))
}
