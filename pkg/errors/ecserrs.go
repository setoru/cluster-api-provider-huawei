package errors

import (
	"net/http"
)

const (
	ErrECSNotFound = "Ecs.0114"
)

// ECSErrorHandler handles ECS-specific errors.
type ECSErrorHandler struct {
	BaseErrorHandler
}

func (e ECSErrorHandler) IsNotFound(err error) bool {
	return e.StatusCode(err) == http.StatusNotFound && e.ErrorCode(err) == ErrECSNotFound
}
