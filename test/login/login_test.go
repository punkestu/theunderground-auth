package login

import (
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity/mocks"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"github.com/punkestu/theunderground-auth/internal/entity/response"
	"github.com/punkestu/theunderground-auth/internal/handler"
	mocks2 "github.com/punkestu/theunderground-auth/internal/lib/mocks"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"github.com/punkestu/theunderground-auth/test/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const endPoint = "/login"

var app *fiber.App
var e *mocks.Entity
var jwt *mocks2.Jwt
var dummyUser1 = &object.User{
	ID:       "1234",
	Username: "minerva",
	Email:    "test@mail.com",
	Password: "test1234",
	Key:      "user1234",
}
var dummyToken = "token1234"
var IdentifierNotFound = response.Error{
	Field:   "Identifier",
	Message: "Identifier not found",
}
var PasswordIsWrong = response.Error{
	Field:   "Password",
	Message: "Password is wrong",
}

func TestLogin(t *testing.T) {
	app = fiber.New()
	jwt = mocks2.NewJwt(t)
	e = mocks.NewEntity(t)
	u := usecase.NewUserUsecase(e)
	h := handler.NewUserHandler(u, jwt)
	app.Post(endPoint, h.Login)
	t.Run("Success Using Email", LoginWithEmailSuccess)
	t.Run("Success Using Username", LoginWithUsernameSuccess)
	t.Run("Failed User Not Found", LoginUserNotFoundFailed)
	t.Run("Failed Wrong Password", LoginWrongPasswordFailed)
}

func LoginWrongPasswordFailed(t *testing.T) {
	e.On("Login", "minerva", "test123").Return("", object.OneError(http.StatusUnauthorized, PasswordIsWrong.GenError()))
	req, err := util.SendRequest(http.MethodPost, endPoint, request.Login{
		Identifier: "minerva",
		Password:   "test123",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)

	assert.Len(t, resBody.Errors, 1)
	assert.Equal(t, resBody.Errors[0].Field, "Password")
}

func LoginUserNotFoundFailed(t *testing.T) {
	e.On("Login", "min", "test1234").Return("", object.OneError(http.StatusUnauthorized, IdentifierNotFound.GenError()))
	req, err := util.SendRequest(http.MethodPost, endPoint, request.Login{
		Identifier: "min",
		Password:   "test1234",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)

	assert.Len(t, resBody.Errors, 1)
	assert.Equal(t, resBody.Errors[0].Field, IdentifierNotFound.Field)
	assert.Equal(t, resBody.Errors[0].Message, IdentifierNotFound.Message)
}

func LoginWithEmailSuccess(t *testing.T) {
	e.On("Login", "test@mail.com", "test1234").Return(
		dummyUser1.ID+"|"+dummyUser1.Key,
		object.NoError(),
	)
	jwt.On("Sign", dummyUser1.ID+"|"+dummyUser1.Key).Return(dummyToken)
	req, err := util.SendRequest(http.MethodPost, endPoint, request.Login{
		Identifier: "test@mail.com",
		Password:   "test1234",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resBody response.AuthSuccess
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
	assert.Equal(t, dummyToken, resBody.Token)
}

func LoginWithUsernameSuccess(t *testing.T) {
	e.On("Login", "minerva", "test1234").Return(
		dummyUser1.ID+"|"+dummyUser1.Key,
		object.NoError(),
	)
	jwt.On("Sign", dummyUser1.ID+"|"+dummyUser1.Key).Return(dummyToken)
	req, err := util.SendRequest(http.MethodPost, endPoint, request.Login{
		Identifier: "minerva",
		Password:   "test1234",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resBody response.AuthSuccess
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
	assert.Equal(t, dummyToken, resBody.Token)
}
