package service

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"tongla-account/entity"
	"tongla-account/util"
)

type AddServiceRequest struct {
	RedirectUri string `json:"redirect_uri" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Issuer      string `json:"issuer" validate:"required"`
}

func (o openIdService) HandleAddServiceRouter(c *fiber.Ctx) error {
	var register AddServiceRequest

	err := util.ValidateRequest(c, &register)
	if err != nil {
		panic(err)
	}

	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	if !strings.HasPrefix(register.RedirectUri, "http://") && !strings.HasPrefix(register.RedirectUri, "https://") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "redirectUri must start with http:// or https://",
		})
	}

	service, err := o.serviceRepository.CreateService(&entity.Service{
		RedirectUri:  register.RedirectUri,
		Name:         register.Name,
		Issuer:       register.Issuer,
		Owner:        user.ID,
		Scope:        "openid,email,profile",
		GrantType:    "authorization_code",
		ResponseType: "code",
	})

	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(service.ToResponse())
}
