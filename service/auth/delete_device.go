package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
	"tongla-account/util"
)

type DeleteDeviceRequest struct {
	SessionId string `json:"session_id" validate:"required"`
}

func (a authService) HandleDeleteDeviceRouter(c *fiber.Ctx) error {
	var request DeleteDeviceRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	token, err := a.jsonWebTokenRepository.GetTokenById(request.SessionId)
	if err != nil {
		panic(err)
	}

	if token.AccountId != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to delete this device",
		})
	}

	err = a.jsonWebTokenRepository.RevokedAllActiveTokenByRefId(request.SessionId)
	if err != nil {
		panic(err)
	}

	_ = a.notificationRepository.SendNotification(&entity.Notification{
		Type:    entity.NotificationWeb,
		Email:   user.Email,
		Title:   "Device Deleted",
		Content: fmt.Sprintf(util.GetWebNotificationContent("deviceDelete"), token.DeviceID, token.Issuer),
		Reason:  "alert",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
