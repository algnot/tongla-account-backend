package util

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateRequest(c *fiber.Ctx, entity any) error {
	body := c.Body()

	if err := json.Unmarshal(body, entity); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(entity); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
