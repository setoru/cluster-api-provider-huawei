package errors

import (
	"net/http"
)

const (
	ErrVPCNotFound = "VPC.0202"
)

// VPCErrorHandler handles VPC-specific errors.
type VPCErrorHandler struct {
	BaseErrorHandler
}

func (v VPCErrorHandler) IsNotFound(err error) bool {
	statusCode, errorCode := v.StatusCode(err), v.ErrorCode(err)
	return (statusCode == http.StatusNotFound || statusCode == http.StatusInternalServerError) && errorCode == ErrVPCNotFound
}
