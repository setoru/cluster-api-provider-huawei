package errors

import (
	"net/http"
)

const (
	ErrAPINotExist     = "APIGW.0101"
	ErrELBDataConflict = "ELB.8907"
)

// ELBErrorHandler handles ELB-specific errors.
type ELBErrorHandler struct {
	BaseErrorHandler
}

func (e ELBErrorHandler) IsNotFound(err error) bool {
	return e.StatusCode(err) == http.StatusNotFound && e.ErrorCode(err) == ErrAPINotExist
}

func (e ELBErrorHandler) IsExists(err error) bool {
	return e.StatusCode(err) == http.StatusConflict && e.ErrorCode(err) == ErrELBDataConflict
}
