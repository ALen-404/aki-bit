package errcode

import (
	"errors"
	"net/http"
)

func CodeFromError(err error) int {
	switch {
	case errors.Is(err, ErrCmcApiError):
		return http.StatusBadRequest
	case errors.Is(err, ErrHttpError):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrCmcNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
