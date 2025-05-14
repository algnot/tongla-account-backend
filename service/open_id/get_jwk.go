package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/gofiber/fiber/v2"
	"math/big"
	"os"
)

func (o openIdService) HandleJWKSRouter(c *fiber.Ctx) error {
	pubKey, err := loadRSAPublicKey("./secret/public.pem")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load public key",
			"detail": err.Error(),
		})
	}

	jwk := map[string]interface{}{
		"keys": []map[string]string{
			{
				"kty": "RSA",
				"use": "sig",
				"alg": "RS256",
				"kid": "tongla.dev",
				"n":   base64Url(pubKey.N),
				"e":   base64Url(big.NewInt(int64(pubKey.E))),
			},
		},
	}

	return c.JSON(jwk)
}

func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func base64Url(n *big.Int) string {
	return base64.RawURLEncoding.EncodeToString(n.Bytes())
}
