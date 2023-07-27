package request

type IdType int

const (
	KeyType             IdType = 1
	UsernameOrEmailType        = 2
)

type Login struct {
	IdentifierType IdType
	Identifier     string `json:"identifier"`
	Password       string `json:"password"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
