package lib

import "time"

type CustomError struct {
	Provider  int         `json:"provider"`
	Error     interface{} `json:"error"`
	CreatedAt time.Time   `json:"created_at"`
}

func CreateError(provider int, err interface{}) *CustomError {
	return &CustomError{
		Provider:  provider,
		Error:     err,
		CreatedAt: time.Now(),
	}
}
