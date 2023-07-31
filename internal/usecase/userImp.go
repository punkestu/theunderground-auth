package usecase

import (
	"errors"
	"github.com/punkestu/theunderground-auth/internal/entity"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"net/http"
)

type UserUsecase struct {
	userEntity entity.Entity
}

func NewUserUsecase(userEntity entity.Entity) *UserUsecase {
	return &UserUsecase{userEntity: userEntity}
}

func (u UserUsecase) Login(r request.Login) (string, object.Error) {
	if r.IdentifierType == request.KeyType {
		return u.userEntity.LoginWithKey(r.Identifier)
	} else if r.IdentifierType == request.UsernameOrEmailType {
		return u.userEntity.Login(r.Identifier, r.Password)
	} else {
		return "", object.OneError(http.StatusBadRequest, errors.New("identifier not valid"))
	}
}

func (u UserUsecase) Register(r request.Register) (string, object.Error) {
	return u.userEntity.Register(r)
}

func (u UserUsecase) GetUser(ID string) (*object.User, object.Error) {
	return u.userEntity.GetUser(ID)
}
