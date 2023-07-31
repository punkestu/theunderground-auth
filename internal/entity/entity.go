package entity

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/entity/request"
	"github.com/punkestu/theunderground-auth/internal/repo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type Entity interface {
	Login(string, string) (string, object.Error)
	LoginWithKey(string) (string, object.Error)
	Register(request.Register) (string, object.Error)
	GetUser(string) (*object.User, object.Error)
}

type UserEntity struct {
	repo.Repo
}

func NewUserEntity(repo repo.Repo) *UserEntity {
	return &UserEntity{Repo: repo}
}

func (e *UserEntity) Login(identifier, password string) (string, object.Error) {
	user, err := e.Repo.GetByUsernameOrEmail(identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", object.OneError(http.StatusUnauthorized, errors.New("Identifier:Identifier not found"))
		}
		return "", object.OneError(http.StatusInternalServerError, errors.New("server error"))
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", object.OneError(http.StatusUnauthorized, errors.New("Password:Password is wrong"))
	}
	return user.ID + "|" + user.Key, object.NoError()
}

func (e *UserEntity) LoginWithKey(key string) (string, object.Error) {
	user, err := e.Repo.GetByKey(key)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", object.OneError(http.StatusUnauthorized, errors.New("Identifier:Key not found"))
		}
		return "", object.OneError(http.StatusInternalServerError, errors.New("server error"))
	}
	return user.ID + "|" + user.Key, object.NoError()
}

func (e *UserEntity) Register(r request.Register) (string, object.Error) {
	id := uuid.New()
	key := strings.Split(uuid.New().String(), "-")[0]
	password, err := bcrypt.GenerateFromPassword([]byte(r.Password), 10)
	if err != nil {
		return "", object.OneError(http.StatusInternalServerError, errors.New("server error"))
	}
	err = e.Repo.Create(object.User{
		ID:       id.String(),
		Username: r.Username,
		Email:    r.Email,
		Password: string(password),
		Key:      key,
	})
	if err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) {
			if mySqlErr.Number == 1062 {
				return "", object.OneError(http.StatusUnauthorized, errors.New("Identifier:Identifier is used"))
			}
		}
		return "", object.OneError(http.StatusInternalServerError, errors.New("server error"))
	}
	return id.String() + "|" + key, object.NoError()
}

func (e *UserEntity) GetUser(ID string) (*object.User, object.Error) {
	user, err := e.Repo.GetByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, object.OneError(http.StatusNotFound, errors.New("ID:ID not found"))
		}
		log.Println(err)
		return nil, object.OneError(http.StatusInternalServerError, errors.New("server error"))
	}
	return &user, object.NoError()
}
