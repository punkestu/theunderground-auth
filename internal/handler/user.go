package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"github.com/punkestu/theunderground-auth/internal/entity/response"
	"github.com/punkestu/theunderground-auth/internal/lib"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"net/http"
	"strings"
)

type User struct {
	user usecase.User
	jwt  lib.Jwt
}

func NewUserHandler(user usecase.User, jwt lib.Jwt) *User {
	return &User{user: user, jwt: jwt}
}

func (u User) Login(c *fiber.Ctx) error {
	var r request.Login
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrors(
			object.OneError(http.StatusInternalServerError, err),
		))
	}
	r.IdentifierType = request.UsernameOrEmailType
	res, err := u.user.Login(r)
	if err.IsError() {
		return c.Status(err.Status).JSON(response.NewErrors(err))
	}
	return c.JSON(response.AuthSuccess{Token: u.jwt.Sign(res)})
}

func (u User) LoginWithKey(c *fiber.Ctx) error {
	f, fErr := c.FormFile("credential")
	if fErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrors(
			object.OneError(http.StatusInternalServerError, fErr),
		))
	}
	buffer, fErr := f.Open()
	if fErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrors(
			object.OneError(http.StatusInternalServerError, fErr),
		))
	}
	token := make([]byte, 256)
	_, fErr = buffer.Read(token)
	if fErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrors(
			object.OneError(http.StatusInternalServerError, fErr),
		))
	}
	res, err := u.user.Login(request.Login{
		IdentifierType: request.KeyType,
		Identifier:     strings.TrimRight(string(token), "\x00"),
	})
	if err.IsError() {
		return c.Status(err.Status).JSON(response.NewErrors(err))
	}
	return c.JSON(response.AuthSuccess{Token: u.jwt.Sign(res)})
}

func (u User) Register(c *fiber.Ctx) error {
	var r request.Register
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrors(
			object.OneError(http.StatusInternalServerError, err),
		))
	}
	res, err := u.user.Register(r)
	if err.IsError() {
		return c.Status(err.Status).JSON(response.NewErrors(err))
	}
	return c.JSON(response.AuthSuccess{Token: u.jwt.Sign(res)})
}

func (u User) GetUser(c *fiber.Ctx) error {
	var r request.GetUser
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrors(
			object.OneError(http.StatusInternalServerError, err),
		))
	}
	res := u.jwt.Parse(r.Token)
	if res == "" {
		return c.Status(http.StatusForbidden).JSON(response.NewErrors(object.OneError(http.StatusForbidden, errors.New("Token:Token invalid"))))
	}
	id := strings.Split(res, "|")[0]
	user, err := u.user.GetUser(id)
	if err.IsError() {
		return c.Status(err.Status).JSON(response.NewErrors(err))
	}
	return c.JSON(response.UserFiltered{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Key:      user.Key,
		Pub:      user.Key, // TODO encrypt Pub and Priv
		Priv:     user.Key,
	})
}
