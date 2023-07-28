package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"github.com/punkestu/theunderground-auth/internal/entity/response"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"net/http"
	"strings"
)

type User struct {
	user usecase.User
}

func NewUserHandler(user usecase.User) *User {
	return &User{user: user}
}

func (u User) Login(c *fiber.Ctx) error {
	var r request.Login
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	r.IdentifierType = request.UsernameOrEmailType
	res, err := u.user.Login(r)
	if err.IsError() {
		return c.Status(err.Status).JSON(response.NewErrors(err))
	}
	// add token generator here
	return c.JSON(response.AuthSuccess{Token: res})
}

func (u User) LoginWithKey(c *fiber.Ctx) error {
	f, fErr := c.FormFile("credential")
	if fErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fErr)
	}
	buffer, fErr := f.Open()
	if fErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fErr)
	}
	token := make([]byte, 256)
	_, fErr = buffer.Read(token)
	if fErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fErr)
	}
	res, err := u.user.Login(request.Login{
		IdentifierType: request.KeyType,
		Identifier:     strings.TrimRight(string(token), "\x00"),
	})
	if err.IsError() {
		return c.Status(err.Status).JSON(response.NewErrors(err))
	}
	return c.JSON(response.AuthSuccess{Token: res})
}

func (u User) Register(c *fiber.Ctx) error {
	var r request.Register
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	res, err := u.user.Register(r)
	if err.IsError() {
		return c.Status(err.Status).JSON(err.Errors)
	}
	// add token generator here
	return c.JSON(response.AuthSuccess{Token: res})
}
