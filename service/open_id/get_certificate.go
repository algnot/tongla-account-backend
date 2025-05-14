package service

import "github.com/gofiber/fiber/v2"

func (o openIdService) HandleCertificateRouter(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"issuer":                                o.config.ServerConfig.BackendPath,
		"authorization_endpoint":                o.config.ServerConfig.FrontendPath + "/authorize",
		"token_endpoint":                        o.config.ServerConfig.BackendPath + "/openid/token",
		"userinfo_endpoint":                     o.config.ServerConfig.BackendPath + "/openid/userinfo",
		"jwks_uri":                              o.config.ServerConfig.BackendPath + "/openid/.well-known/jwks.json",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	})
}
