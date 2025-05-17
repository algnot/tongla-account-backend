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
		panic(err)
	}

	key, err := a.apiKeysRepository.RotateKeyByName(apiKeysRequest.Name)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(key)
}
