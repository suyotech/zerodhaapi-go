package kiteconnect

import "fmt"

type APIError struct {
	StatusCode int
	Status     string `json:"status"`
	Message    string `json:"message"`
	ErrorType  string `json:"error_type"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.ErrorType == "" {
		return fmt.Sprintf("kiteconnect: %s", e.Message)
	}
	return fmt.Sprintf("kiteconnect: %s: %s", e.ErrorType, e.Message)
}

func IsTokenError(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.ErrorType == "TokenException"
}
