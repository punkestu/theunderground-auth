package object

type Error struct {
	Status int     `json:"status"`
	Errors []error `json:"errors"`
}

func (e Error) IsError() bool {
	if e.Status == NoError().Status {
		return false
	} else {
		return true
	}
}

func NoError() Error {
	return Error{
		Status: 0,
	}
}

func OneError(status int, theError error) Error {
	return Error{
		Status: status,
		Errors: []error{
			theError,
		},
	}
}
