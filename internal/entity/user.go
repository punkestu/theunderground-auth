package entity

type User struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Key      string `json:"key"`
}
