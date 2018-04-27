package handler

import (
	"net/http"

	"github.com/src-d/gitbase-playground/server/serializer"
)

// Parse : placeholder method
func Parse() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}

// Filter : placeholder method
func Filter() RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return nil, serializer.NewHTTPError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	}
}
