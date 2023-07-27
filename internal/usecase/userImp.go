package usecase

import (
	"errors"
	"github.com/punkestu/theunderground-auth/internal/entity"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"github.com/punkestu/theunderground-auth/internal/repo"
	"net/http"
)

type UserUsecase struct {
	repo.Repo
}

func NewUserUsecase(repo repo.Repo) *UserUsecase {
	return &UserUsecase{Repo: repo}
}

func (u UserUsecase) Login(r request.Login) (string, entity.Error) {
	if r.IdentifierType == request.KeyType {
		user, err := u.GetByKey(r.Identifier)
		if err.IsError() {
			return "", entity.NoError()
		}
		return user.ID + user.Key, entity.NoError()
	} else if r.IdentifierType == request.UsernameOrEmailType {
		user, err := u.GetByUsernameOrEmail(r.Identifier)
		if err.IsError() {
			return "", err
		}
		if user.Password != r.Password {
			return "", entity.OneError(http.StatusNotFound, errors.New("Password:Password is wrong"))
		}
		return user.ID + "-" + user.Key, entity.NoError()
	} else {
		return "", entity.OneError(http.StatusBadRequest, errors.New("identifier not valid"))
	}
}

func (u UserUsecase) Register(r request.Register) (string, entity.Error) {
	id, err := u.Create(entity.User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	})
	if err.IsError() {
		return "", err
	}
	user, err := u.GetByID(id)
	if err.IsError() {
		return "", err
	}
	return user.ID + "-" + user.Key, entity.NoError()
}
