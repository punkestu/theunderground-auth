package repo

import (
	"github.com/punkestu/theunderground-auth/internal/entity/object"
)

type Repo interface {
	GetByID(string) (object.User, error)
	GetByUsernameOrEmail(string) (object.User, error)
	GetByKey(string) (object.User, error)
	Create(object.User) error
}
