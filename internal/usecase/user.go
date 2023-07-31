package usecase

import (
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
)

type User interface {
	Login(request.Login) (string, object.Error)
	Register(request.Register) (string, object.Error)
	GetUser(string) (*object.User, object.Error)
}
