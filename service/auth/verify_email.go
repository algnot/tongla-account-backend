package service

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"time"
	"tongla-account/entity"
	"tongla-account/util"
)

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

func (a authService) HandleVerifyEmailRouter(c *fiber.Ctx) error {
	var request VerifyEmailRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	token, err := a.tokenRepository.FindKeyByToken(request.Token)
	if err != nil {
		panic(err)
	}

	if token.Type != entity.TokenVerifyEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token Type is not VerifyEmail",
		})
	}

	if token.ExpireAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Link is expired try to login",
		})
	}

	user, err := a.accountRepository.FindById(token.AccountID)
	if err != nil {
		panic(err)
	}

	if user.IsVerified {
		panic(err)
	}

	secret, err := a.accountRepository.GenerateSecret(user)
	if err != nil {
		panic(err)
	}

	qrBytes, err := qrcode.Encode(secret, qrcode.Medium, 256)
	if err != nil {
		panic(err)
	}

	base64QR := base64.StdEncoding.EncodeToString(qrBytes)
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64QR)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"qr_code": dataURL,
	})
}
