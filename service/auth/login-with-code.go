package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"tongla-account/util"
)

type LoginWithCodeRequest struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

func (a authService) HandleLoginWithCodeRouter(c *fiber.Ctx) error {
	var request LoginWithCodeRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := a.accountRepository.FindByEmail(request.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !user.IsVerified {
		err = a.accountRepository.SendVerifyEmail(user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"redirect": a.config.ServerConfig.FrontendPath + "/sign-up/verify?email=" + request.Email,
		})
	}

	userSecret := a.encryptorRepository.Decrypt(user.Secret)
	valid := totp.Validate(request.Code, userSecret)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is invalid",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user.ToResponse(a.encryptorRepository.Decrypt))
}
