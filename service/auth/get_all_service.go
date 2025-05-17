package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

type GetAllServiceResponse struct {
	Services []entity.ServicesResponse `json:"services"`
}

func (a authService) HandleGetAllServiceRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	services, err := a.serviceRepository.GetAllServiceByAccountId(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get service by account",
		})
	}

	var response []entity.ServicesResponse
	for _, s := range services {
		response = append(response, s.ToResponse())
	}

	return c.Status(fiber.StatusOK).JSON(GetAllServiceResponse{
		Services: response,
	})
}
