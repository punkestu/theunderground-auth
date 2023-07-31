package response

import (
	"errors"
	"github.com/punkestu/theunderground-auth/internal/entity/object"
	"strings"
)

type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (e Error) GenError() error {
	return errors.New(e.Field + ":" + e.Message)
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func NewErrors(theError object.Error) Errors {
	var theErrors Errors
	for _, err := range theError.Errors {
		e := strings.Split(err.Error(), ":")
		if len(e) > 1 {
			theErrors.Errors = append(theErrors.Errors, Error{
				Field:   e[0],
				Message: e[1],
			})
		} else {
			theErrors.Errors = append(theErrors.Errors, Error{
				Field:   "",
				Message: e[0],
			})
		}
	}
	return theErrors
}
