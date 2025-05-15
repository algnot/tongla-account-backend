package service

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
	"tongla-account/entity"
)

func (o openIdService) HandleGetTokenRouter(c *fiber.Ctx) error {
	clientID := ""
	clientSecret := ""
	code := c.FormValue("code")
	redirectURI := c.FormValue("redirect_uri")
	grantType := c.FormValue("grant_type")

	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Basic ") {
		raw, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		parts := strings.SplitN(string(raw), ":", 2)
		clientID = parts[0]
		clientSecret = parts[1]
	} else {
		clientID = c.FormValue("client_id")
		clientSecret = c.FormValue("client_secret")
	}

	if grantType != "authorization_code" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid grant type",
		})
	}

	client, err := o.serviceRepository.GetByClientId(clientID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if client.ClientSecret != clientSecret {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid client secret",
		})
	}

	if client.RedirectUri != redirectURI {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid redirect uri",
		})
	}

	token, err := o.tokenRepository.FindKeyByToken(code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if token.Type != entity.TokenAuthCode {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid token type",
		})
	}

	if token.Ref != clientID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid client id",
		})
	}

	if token.ExpireAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token is expired",
		})
	}

	if token.Used {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is already used",
		})
	}

	user, err := o.accountRepository.FindById(token.AccountID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token.Used = true
	_, err = o.tokenRepository.UpdateToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	activeRefreshToken, _ := o.jsonWebTokenRepository.GetActiveRefreshTokenByClientId(clientID)

	var accessToken string
	if activeRefreshToken == nil {
		jwtToken, err := o.jsonWebTokenRepository.GenerateToken(user, client.Issuer, client.Issuer, client.ClientId, client.Name, client.ClientId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		accessToken = jwtToken.AccessToken
	} else {
		accessTokenEnt, err := o.jsonWebTokenRepository.GenerateAccessToken(user, client.Issuer, client.Issuer, client.ClientId, client.Name, client.ClientId, activeRefreshToken.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		accessToken = accessTokenEnt
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   60 * 60 * 10,
		"id_token":     accessToken,
		//"refresh_token": jwtToken.RefreshToken,
	})
}
