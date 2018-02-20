package accessor

import (
	"errors"
	"net/http"
)

var ErrForbidden = errors.New("Forbidden")
var ErrNotAuthorized = errors.New("Not Authorized")
var ErrNotFound = errors.New("Not Found")

func HttpStatus(err error) int {
	switch err {
	case ErrNotAuthorized:
		return http.StatusUnauthorized
	case ErrNotFound:
		return http.StatusNotFound
	case ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
