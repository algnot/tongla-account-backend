package service

import (
	"github.com/gofiber/fiber/v2"
	"tongla-account/entity"
	"tongla-account/util"
)

type RegisterRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
}

func (a authService) HandleRegisterRouter(c *fiber.Ctx) error {
	var registerRequest RegisterRequest

	err := util.ValidateRequest(c, &registerRequest)
	if err != nil {
		panic(err)
	}

	createdAccount, err := a.accountRepository.CreateAccount(&entity.Account{
		Username:  a.encryptorRepository.Encrypt(registerRequest.Username),
		Email:     a.encryptorRepository.Encrypt(registerRequest.Email),
		Firstname: a.encryptorRepository.Encrypt(registerRequest.Firstname),
		Lastname:  a.encryptorRepository.Encrypt(registerRequest.Lastname),
		Birthdate: nil,
	})

	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusCreated).JSON(createdAccount.ToResponse(a.encryptorRepository.Decrypt))
}
