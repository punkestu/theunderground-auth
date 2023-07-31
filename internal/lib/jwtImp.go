package lib

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWT struct{}

func NewJWT() Jwt {
	return &JWT{}
}

func (j *JWT) Sign(payload string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(1 * time.Hour)},
		Subject:   payload,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return ""
	}
	return tokenString
}

func (j *JWT) Parse(token string) string {
	if len(token) < 8 {
		return ""
	}
	res, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if res == nil || err != nil {
		return ""
	}
	claims, ok := res.Claims.(jwt.MapClaims)
	if !ok || !res.Valid {
		return ""
	}
	sub, _ := claims.GetSubject()
	return sub
}
