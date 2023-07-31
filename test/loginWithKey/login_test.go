package loginWithKey

import (
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity/mocks"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/entity/response"
	"github.com/punkestu/theunderground-auth/internal/handler"
	mocks2 "github.com/punkestu/theunderground-auth/internal/lib/mocks"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"github.com/punkestu/theunderground-auth/test/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const endPoint = "/key"

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

func TestLoginWithKey(t *testing.T) {
	app = fiber.New()
	jwt = mocks2.NewJwt(t)
	e = mocks.NewEntity(t)
	u := usecase.NewUserUsecase(e)
	h := handler.NewUserHandler(u, jwt)
	app.Post(endPoint, h.LoginWithKey)
	t.Run("Success", Success)
	t.Run("Failed Key not found", KeyNotFoundFailed)
}

func KeyNotFoundFailed(t *testing.T) {
	e.On("LoginWithKey", "user123").Return("", object.OneError(http.StatusUnauthorized, IdentifierNotFound.GenError()))
	req, err := util.SendFileRequest(http.MethodPost, endPoint, "wrongCred.key", nil)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
}

func Success(t *testing.T) {
	e.On("LoginWithKey", "user1234").Return(dummyUser1.ID+"|"+dummyUser1.Key, object.NoError())
	req, err := util.SendFileRequest(http.MethodPost, endPoint, "credential.key", nil)
	jwt.On("Sign", dummyUser1.ID+"|"+dummyUser1.Key).Return(dummyToken)

	res, err := app.Test(req)
	assert.Nil(t, err)
	var resBody response.AuthSuccess
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
	assert.Equal(t, dummyToken, resBody.Token)
}
