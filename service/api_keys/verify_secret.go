package service

import (
	"github.com/gofiber/fiber/v2"
)

func (a apiKeysService) HandleVerifyGetRouter(c *fiber.Ctx) error {
	apiKey := c.Get("X-Api-Key")

	result, err := a.apiKeysRepository.VerifyKey(apiKey)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
