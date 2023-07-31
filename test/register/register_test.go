package register

import (
	"errors"
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

const endPoint = "/register"

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

func TestRegister(t *testing.T) {
	app = fiber.New()
	jwt = mocks2.NewJwt(t)
	e = mocks.NewEntity(t)
	u := usecase.NewUserUsecase(e)
	h := handler.NewUserHandler(u, jwt)
	app.Post(endPoint, h.Register)
	t.Run("Failed Username existed", UserNameExistedFailed)
	t.Run("Success", Success)
}

func UserNameExistedFailed(t *testing.T) {
	e.On("Register", request.Register{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: "test",
	}).Return("", object.OneError(http.StatusUnauthorized, errors.New("Username:Username is used")))

	req, err := util.SendRequest(http.MethodPost, endPoint, request.Register{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: "test",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
}

func Success(t *testing.T) {
	e.On("Register", request.Register{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: dummyUser1.Password,
	}).Return(dummyUser1.ID+"|"+dummyUser1.Key, object.NoError())
	jwt.On("Sign", dummyUser1.ID+"|"+dummyUser1.Key).Return(dummyToken)
	req, err := util.SendRequest(http.MethodPost, endPoint, request.Register{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: dummyUser1.Password,
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
