package lib

type Jwt interface {
	Sign(payload string) string
	Parse(token string) string
}
