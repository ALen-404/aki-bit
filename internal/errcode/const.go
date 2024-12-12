package errcode

import "errors"

var (
	ErrCmcApiError = errors.New("failed to retrieve content from cmc")
	ErrHttpError   = errors.New("http request occur error")
	ErrCmcNotFound = errors.New("resource not found")
)
