package response

type AuthSuccess struct {
	Token string `json:"token"`
}

type UserFiltered struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Key      string `json:"key"`
	Pub      string `json:"pub"`
	Priv     string `json:"priv"`
}
