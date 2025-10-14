package http

import (
	"errors"

	"qrcodegen/internal/dto"
	"qrcodegen/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	validate    *validator.Validate
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(validate *validator.Validate, userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		validate:    validate,
		userUseCase: userUseCase,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body      dto.RegisterRequest  true  "User registration data"
// @Success 201   {object}  dto.RegisterResponse
// @Failure 400   {object}  dto.GenericError
// @Failure 409   {object}  dto.GenericError
// @Failure 500   {object}  dto.GenericError
// @Router /register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.userUseCase.Register(c.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Login godoc
// @Summary Log in a user
// @Description Log in a user with email and password, returns a JWT token in a cookie
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body      dto.LoginRequest  true  "User login data"
// @Success 200   {string}  string "OK"
// @Failure 400   {object}  dto.GenericError
// @Failure 401   {object}  dto.GenericError
// @Failure 500   {object}  dto.GenericError
// @Router /login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, expirationTime, err := h.userUseCase.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  expirationTime,
		HTTPOnly: true,
	})

	return c.SendStatus(fiber.StatusOK)
}
