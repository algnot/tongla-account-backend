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
		panic(err)
	}

	user, err := a.accountRepository.FindByEmail(request.Email)
	if err != nil {
		panic(err)
	}

	if !user.IsVerified {
		err = a.accountRepository.SendVerifyEmail(user)
		if err != nil {
			panic(err)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"redirect": a.config.ServerConfig.FrontendPath + "/sign-up/verify?email=" + request.Email,
		})
	}

	env := a.config.CommonConfig.Env
	userSecret := a.encryptorRepository.Decrypt(user.Secret)
	valid := totp.Validate(request.Code, userSecret)
	if !valid && env != "local" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is invalid",
		})
	}

	userAgent := c.Get("User-Agent")
	deviceID := c.Get("Device-ID")
	token, err := a.jsonWebTokenRepository.GenerateToken(user, "tongla.dev", "tongla.dev", userAgent, deviceID, "")
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(token)
}
