package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/util"
)

type LoginWithEmailRequest struct {
	Email string `json:"email" validate:"required"`
}

func (a authService) HandleRequestLoginWithEmailRouter(c *fiber.Ctx) error {
	var request LoginWithEmailRequest

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

	err = a.accountRepository.SendLoginLinkWithEmail(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
