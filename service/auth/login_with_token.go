package service

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"tongla-account/entity"
	"tongla-account/util"
)

type LoginWithTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

func (a authService) HandleLoginWithTokenRouter(c *fiber.Ctx) error {
	var request LoginWithTokenRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := a.tokenRepository.FindKeyByToken(request.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if token.Type != entity.TokenLogin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token Type is not login",
		})
	}

	if token.ExpireAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Link is expired",
		})
	}

	if token.Used {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is already used",
		})
	}

	user, err := a.accountRepository.FindById(token.AccountID)
	if err != nil {
		panic(err)
	}

	token.Used = true
	_, err = a.tokenRepository.UpdateToken(token)
	if err != nil {
		panic(err)
	}

	userAgent := c.Get("User-Agent")
	deviceID := c.Get("Device-ID")
	jwtToken, err := a.jsonWebTokenRepository.GenerateToken(user, "tongla.dev", "tongla.dev", userAgent, deviceID, "")
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(jwtToken)
}
