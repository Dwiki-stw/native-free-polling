package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func (e *AppError) WriteError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	switch e.Code {
	case "INVALID_INPUT":
		status = http.StatusBadRequest
	case "EMAIL_EXIST":
		status = http.StatusConflict
	case "AUTH_FAILED":
		status = http.StatusUnauthorized
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    e.Code,
		"message": e.Message,
	})
}
