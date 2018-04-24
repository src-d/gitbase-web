package handler

import (
	"net/http"

	"github.com/src-d/gitbase-playground/server/serializer"
)

// Version returns a function that returns a *serializer.Response
// with a current version of server
func Version(version string) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		return serializer.NewVersionResponse(version), nil
	}
}
