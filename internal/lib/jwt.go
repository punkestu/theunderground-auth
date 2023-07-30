package lib

type Jwt interface {
	Sign(payload any) string
	Validate(token string) bool
}
