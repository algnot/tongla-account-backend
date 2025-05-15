package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
)

type ServicesResponse struct {
	Name                   string `json:"name"`
	ClientId               string `json:"client_id"`
	ClientSecret           string `json:"client_secret"`
	RedirectUri            string `json:"redirect_uri"`
	Issuer                 string `json:"issuer"`
	Scopes                 string `json:"scopes"`
	GrantTypes             string `json:"grant_types"`
	ResponseType           string `json:"response_type"`
	OpenidConfigurationUri string `json:"openid_configuration_uri"`
}

type GetAllServiceResponse struct {
	Services []ServicesResponse `json:"services"`
}

func (a authService) HandleGetAllServiceRouter(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	services, err := a.serviceRepository.GetAllServiceByAccountId(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get service by account",
		})
	}

	var response []ServicesResponse
	for _, s := range services {
		response = append(response, ServicesResponse{
			Name:                   s.Name,
			ClientId:               s.ClientId,
			ClientSecret:           s.ClientSecret,
			RedirectUri:            s.RedirectUri,
			Issuer:                 s.Issuer,
			Scopes:                 s.Scope,
			GrantTypes:             s.GrantType,
			ResponseType:           s.ResponseType,
			OpenidConfigurationUri: a.config.ServerConfig.BackendPath + "/openid/.well-known/configuration",
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetAllServiceResponse{
		Services: response,
	})
}
