package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
	"tongla-account/util"
)

func (a apiKeysService) HandleRotatePostRouter(c *fiber.Ctx) error {
	var apiKeysRequest entity.ApiKeysRequest

	err := util.ValidateRequest(c, &apiKeysRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	key, err := a.apiKeysRepository.RotateKeyByName(apiKeysRequest.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(key)
}
