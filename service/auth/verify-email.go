package service

import (
	"github.com/gofiber/fiber/v2"
	"time"
	"tongla-account/entity"
	"tongla-account/util"
)

func (a authService) HandleVerifyEmailRouter(c *fiber.Ctx) error {
	var request entity.VerifyEmailRequest

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

	if token.Type != entity.TokenVerifyEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token Type is not VerifyEmail",
		})
	}

	if token.Used {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is already used",
		})
	}

	if token.ExpireAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is expired",
		})
	}

	token.Used = true
	_, err = a.tokenRepository.UpdateToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := a.accountRepository.FindById(token.AccountID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User is already verified",
		})
	}

	user.IsVerified = true
	user, err = a.accountRepository.UpdateAccount(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user.ToResponse(a.encryptorRepository.Decrypt))
}
