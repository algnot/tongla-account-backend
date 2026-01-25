package service

import (
	"tongla-account/util"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email string `json:"email" validate:"required"`
}

func (a authService) HandleLoginRouter(c *fiber.Ctx) error {
	var registerRequest LoginRequest

	err := util.ValidateRequest(c, &registerRequest)
	if err != nil {
		panic(err)
	}

	user, err := a.accountRepository.FindByEmail(registerRequest.Email)
	if err != nil {
		panic(err)
	}

	if !user.IsVerified {
		err = a.accountRepository.SendVerifyEmail(user)
		if err != nil {
			panic(err)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"redirect": a.config.ServerConfig.FrontendPath + "/sign-up/verify?email=" + registerRequest.Email,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
