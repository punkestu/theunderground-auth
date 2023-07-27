package usecase

import (
	"github.com/punkestu/theunderground-auth/internal/entity"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
)

type User interface {
	Login(request.Login) (string, entity.Error)
	Register(request.Register) (string, entity.Error)
}
