package response

import (
	"github.com/punkestu/theunderground-auth/internal/entity"
)

type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func NewErrors(theError entity.Error) Errors {
	var theErrors Errors
	for _, err := range theError.Errors {
		theErrors.Errors = append(theErrors.Errors, Error{
			Field:   "",
			Message: err.Error(),
		})
	}
	return theErrors
}
