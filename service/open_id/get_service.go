package service

import (
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
	"tongla-account/entity"
	"tongla-account/util"
)

type GetServiceRequest struct {
	ClientId     string `json:"client_id" validate:"required"`
	RedirectUri  string `json:"redirect_uri" validate:"required"`
	ResponseType string `json:"response_type" validate:"required"`
	Scope        string `json:"scope" validate:"required"`
	Domain       string `json:"domain"`
	State        string `json:"state"`
}

func checkScope(scopes string, check string) bool {
	scopeSet := make(map[string]bool)
	for _, p := range strings.Split(scopes, ",") {
		scopeSet[p] = true
	}

	for _, c := range strings.Fields(check) {
		if !scopeSet[c] {
			return false
		}
	}

	return true
}

func (o openIdService) HandleGetServiceRouter(c *fiber.Ctx) error {
	var request GetServiceRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if request.ResponseType != "code" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid response type",
		})
	}

	client, err := o.serviceRepository.GetByClientId(request.ClientId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if client.RedirectUri != request.RedirectUri {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid redirect_uri",
		})
	}

	if !checkScope(client.Scope, request.Scope) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid scope",
		})
	}

	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	token, err := o.encryptorRepository.GeneratePassphrase(50)
	if err != nil {
		return err
	}
	tokenEnt := &entity.Token{
		AccountID: user.ID,
		Type:      entity.TokenAuthCode,
		Token:     token,
		ExpireAt:  time.Now().Add(30 * time.Minute),
		Ref:       client.ClientId,
	}
	tokenEnt, err = o.tokenRepository.CreateToken(tokenEnt)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"redirect": client.RedirectUri + "?code=" + token + "&state=" + request.State,
	})
}
