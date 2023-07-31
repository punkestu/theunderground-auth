package db

import (
	"database/sql"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"github.com/punkestu/theunderground-auth/internal/repo"
)

type DB struct {
	conn *sql.DB
}

func NewDB(conn *sql.DB) repo.Repo {
	return &DB{conn: conn}
}

func (d *DB) GetByID(id string) (object.User, error) {
	var user object.User
	err := d.conn.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Key)
	return user, err
}
func (d *DB) GetByUsernameOrEmail(identifier string) (object.User, error) {
	var user object.User
	err := d.conn.QueryRow("SELECT * FROM users WHERE email=? OR username=?", identifier, identifier).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Key)
	return user, err
}
func (d *DB) GetByKey(key string) (object.User, error) {
	var user object.User
	err := d.conn.QueryRow("SELECT * FROM users WHERE private_key=?", key).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Key)
	return user, err
}
func (d *DB) Create(user object.User) error {
	_, err := d.conn.Exec("INSERT INTO users VALUES (?, ?, ?, ?, ?)", user.ID, user.Username, user.Email, user.Password, user.Key)
	return err
}
