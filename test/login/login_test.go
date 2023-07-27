package login

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"github.com/punkestu/theunderground-auth/internal/entity/response"
	"github.com/punkestu/theunderground-auth/internal/handler"
	"github.com/punkestu/theunderground-auth/internal/repo/mocks"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"github.com/punkestu/theunderground-auth/test/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const endPoint = "/login"

var app *fiber.App
var r *mocks.Repo
var dummyUser1 = &entity.User{
	ID:       "1234",
	Username: "minerva",
	Email:    "test@mail.com",
	Password: "test1234",
	Key:      "user1234",
}
var UserNotFound = entity.OneError(http.StatusBadRequest, errors.New("user Not Found"))

func TestLogin(t *testing.T) {
	app = fiber.New()
	r = mocks.NewRepo(t)
	u := usecase.NewUserUsecase(r)
	h := handler.NewUserHandler(u)
	app.Post(endPoint, h.Login)
	t.Run("Success Using Email", LoginWithEmailSuccess)
	t.Run("Success Using Username", LoginWithUsernameSuccess)
	t.Run("Failed User Not Found", LoginUserNotFoundFailed)
}

func LoginUserNotFoundFailed(t *testing.T) {
	r.On("GetByUsernameOrEmail", "min").Return(nil, UserNotFound)
	req, err := util.SendRequest(http.MethodPost, endPoint, request.Login{
		Identifier: "min",
		Password:   "test1234",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, &util.GetBodyOptions{Verbose: true})
	assert.Nil(t, err)
}

func LoginWithEmailSuccess(t *testing.T) {
	r.On("GetByUsernameOrEmail", "test@mail.com").Return(
		dummyUser1,
		entity.NoError(),
	)
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
	assert.Equal(t, dummyUser1.ID+"-"+dummyUser1.Key, resBody.Token)
}

func LoginWithUsernameSuccess(t *testing.T) {
	r.On("GetByUsernameOrEmail", "minerva").Return(
		dummyUser1,
		entity.NoError(),
	)
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
	assert.Equal(t, dummyUser1.ID+"-"+dummyUser1.Key, resBody.Token)
}
