package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"time"
	"tongla-account/entity"
	"tongla-account/util"
)

type UpdateUserRequest struct {
	Username  string            `json:"username" validate:"required"`
	Firstname string            `json:"firstname" validate:"required"`
	Lastname  string            `json:"lastname" validate:"required"`
	Gender    entity.GenderType `json:"gender" validate:"required"`
	Birthdate string            `json:"birthdate"`
	Phone     string            `json:"phone"`
	Address   string            `json:"address"`
	Code      string            `json:"code" validate:"required"`
}

func (a authService) HandleUpdateUserRouter(c *fiber.Ctx) error {
	var request UpdateUserRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := c.Locals("user").(*entity.Account)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	userSecret := a.encryptorRepository.Decrypt(user.Secret)
	valid := totp.Validate(request.Code, userSecret)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is invalid",
		})
	}

	var birthdate *time.Time
	if request.Birthdate != "" {
		t, err := time.Parse("2006-01-02", request.Birthdate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid birthdate format (expected YYYY-MM-DD)",
			})
		}
		birthdate = &t
	}

	updatedUser, err := a.accountRepository.UpdateAccount(&entity.Account{
		ID:        user.ID,
		Username:  a.encryptorRepository.Encrypt(request.Username),
		Firstname: a.encryptorRepository.Encrypt(request.Firstname),
		Lastname:  a.encryptorRepository.Encrypt(request.Lastname),
		Phone:     a.encryptorRepository.Encrypt(request.Phone),
		Address:   a.encryptorRepository.Encrypt(request.Address),
		Birthdate: birthdate,
		Gender:    request.Gender,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser.ToResponse(a.encryptorRepository.Decrypt))
}
