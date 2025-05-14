package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

type DeviceResponse struct {
	SessionId string `json:"sessionId"`
	UserAgent string `json:"userAgent"`
	DeviceId  string `json:"deviceId"`
	Issuer    string `json:"issuer"`
	IssuerAt  int64  `json:"issuerAt"`
	Current   bool   `json:"current"`
}

type GetAllDeviceResponse struct {
	Device []DeviceResponse `json:"devices"`
}

func (a authService) HandleGetAllDevice(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	accessToken := c.Locals("token").(*entity.JsonWebToken)
	if accessToken == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token context missing",
		})
	}

	tokens, err := a.jsonWebTokenRepository.GetAllActiveRefreshTokenByAccountId(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tokens",
		})
	}

	var devices []DeviceResponse
	for _, token := range *tokens {
		isCurrent := accessToken.Ref == token.ID
		devices = append(devices, DeviceResponse{
			SessionId: token.ID,
			UserAgent: token.UserAgent,
			DeviceId:  token.DeviceID,
			Issuer:    token.Issuer,
			IssuerAt:  token.Iat,
			Current:   isCurrent,
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetAllDeviceResponse{Device: devices})
}
