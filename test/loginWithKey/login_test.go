package loginWithKey

import (
	"github.com/gofiber/fiber/v2"
	"github.com/punkestu/theunderground-auth/internal/entity"
	"github.com/punkestu/theunderground-auth/internal/entity/response"
	"github.com/punkestu/theunderground-auth/internal/handler"
	"github.com/punkestu/theunderground-auth/internal/repo/mocks"
	"github.com/punkestu/theunderground-auth/internal/usecase"
	"github.com/punkestu/theunderground-auth/test/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const endPoint = "/key"

var app *fiber.App
var r *mocks.Repo
var dummyUser1 = &entity.User{
	ID:       "1234",
	Username: "minerva",
	Email:    "test@mail.com",
	Password: "test1234",
	Key:      "user1234",
}
var IdentifierNotFound = response.Error{
	Field:   "Identifier",
	Message: "Identifier not found",
}

func TestLoginWithKey(t *testing.T) {
	app = fiber.New()
	r = mocks.NewRepo(t)
	u := usecase.NewUserUsecase(r)
	h := handler.NewUserHandler(u)
	app.Post(endPoint, h.LoginWithKey)
	t.Run("Success", Success)
	t.Run("Failed Key not found", KeyNotFoundFailed)
}

func KeyNotFoundFailed(t *testing.T) {
	r.On("GetByKey", "user123").Return(nil, entity.OneError(http.StatusNotFound, IdentifierNotFound.GenError()))
	req, err := util.SendFileRequest(http.MethodPost, endPoint, "wrongCred.key", nil)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
}

func Success(t *testing.T) {
	r.On("GetByKey", "user1234").Return(dummyUser1, entity.NoError())
	req, err := util.SendFileRequest(http.MethodPost, endPoint, "credential.key", nil)

	res, err := app.Test(req)
	assert.Nil(t, err)
	var resBody response.AuthSuccess
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
}
