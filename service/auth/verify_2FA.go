package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"tongla-account/entity"
	"tongla-account/util"
)

type Verify2FARequest struct {
	Token string `json:"token" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

func (a authService) HandleResendVerify2FARouter(c *fiber.Ctx) error {
	var request Verify2FARequest

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

	userSecret := a.encryptorRepository.Decrypt(user.Secret)
	valid := totp.Validate(request.Code, userSecret)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is invalid",
		})
	}

	token.Used = true
	_, err = a.tokenRepository.UpdateToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user.IsVerified = true
	verifyUser, err := a.accountRepository.UpdateAccount(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(verifyUser.ToResponse(a.encryptorRepository.Decrypt))
}
