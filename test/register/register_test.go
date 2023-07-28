package register

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

const endPoint = "/register"

var app *fiber.App
var r *mocks.Repo
var dummyUser1 = &entity.User{
	ID:       "1234",
	Username: "minerva",
	Email:    "test@mail.com",
	Password: "test1234",
	Key:      "user1234",
}

func TestRegister(t *testing.T) {
	app = fiber.New()
	r = mocks.NewRepo(t)
	u := usecase.NewUserUsecase(r)
	h := handler.NewUserHandler(u)
	app.Post(endPoint, h.Register)
	t.Run("Failed Username existed", UserNameExistedFailed)
	t.Run("Success", Success)
}

func UserNameExistedFailed(t *testing.T) {
	r.On("Create", entity.User{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: "test",
	}).Return("", entity.OneError(http.StatusBadRequest, errors.New("Username:Username is used")))

	req, err := util.SendRequest(http.MethodPost, endPoint, request.Register{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: "test",
	}, nil)
	assert.Nil(t, err)

	res, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resBody response.Errors
	err = util.GetBody(res, &resBody, nil)
	assert.Nil(t, err)
}

func Success(t *testing.T) {
	r.On("Create", entity.User{
		Username: dummyUser1.Username,
		Email:    dummyUser1.Email,
		Password: dummyUser1.Password,
	}).Return(dummyUser1.ID, entity.NoError())
	r.On("GetByID", dummyUser1.ID).Return(dummyUser1, entity.NoError())
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
	assert.Equal(t, dummyUser1.ID+"-"+dummyUser1.Key, resBody.Token)
}
