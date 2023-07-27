package repo

import "github.com/punkestu/theunderground-auth/internal/entity"

type Repo interface {
	GetByID(string) (*entity.User, entity.Error)
	GetByUsernameOrEmail(string) (*entity.User, entity.Error)
	GetByKey(string) (*entity.User, entity.Error)
	Create(entity.User) (string, entity.Error)
}
