package errors

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
)

// ErrorHandler defines the interface for handling errors.
type ErrorHandler interface {
	StatusCode(err error) int
	ErrorCode(err error) string
	IsNotFound(err error) bool
}

// BaseErrorHandler provides common error handling methods.
type BaseErrorHandler struct{}

func (b BaseErrorHandler) StatusCode(err error) int {
	if t, ok := err.(*sdkerr.ServiceResponseError); ok {
		return t.StatusCode
	}
	return -1
}

func (b BaseErrorHandler) ErrorCode(err error) string {
	if t, ok := err.(*sdkerr.ServiceResponseError); ok {
		return t.ErrorCode
	}
	return ""
}
