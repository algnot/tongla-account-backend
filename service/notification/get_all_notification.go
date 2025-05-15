package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

type GetAllNotificationsResponse struct {
	Notifications []fiber.Map `json:"notifications"`
}

func (n notificationService) HandleGetAllNotificationsRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	notifications, err := n.notificationRepository.GetNotificationByEmailAndType(user.Email, entity.NotificationWeb)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notifications",
		})
	}

	var response []fiber.Map
	for _, notification := range *notifications {
		response = append(response, notification.ToResponse())
	}

	if response == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"notifications": []fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetAllNotificationsResponse{
		Notifications: response,
	})
}
