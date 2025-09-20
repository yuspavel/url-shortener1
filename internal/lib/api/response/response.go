package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Эта структура будет возвращаться, если возникнет какая-либо ошибка
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

// Анализ ошибки и формирование форматированного ответа (структуры Response)
func ValidateError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s in not a url", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{Status: StatusError, Error: strings.Join(errMsg, ", ")}
}
