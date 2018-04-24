package serializer

import (
	"net/http"
	"strings"
)

// HTTPError defines an Error message as it will be written in the http.Response
type HTTPError interface {
	error
	StatusCode() int
}

// Response encapsulate the content of an http.Response
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
	Errors []HTTPError `json:"errors,omitempty"`
}

type httpError struct {
	Status    int    `json:"status"`
	Title     string `json:"title"`
	Details   string `json:"details,omitempty"`
	MySQLCode uint16 `json:"mysqlCode,omitempty"`
}

// StatusCode returns the Status of the httpError
func (e httpError) StatusCode() int {
	return e.Status
}

// Error returns the string content of the httpError
func (e httpError) Error() string {
	if msg := e.Title; msg != "" {
		return msg
	}

	if msg := http.StatusText(e.Status); msg != "" {
		return msg
	}

	return http.StatusText(http.StatusInternalServerError)
}

// NewHTTPError returns an Error
func NewHTTPError(statusCode int, msg ...string) HTTPError {
	return httpError{Status: statusCode, Title: strings.Join(msg, " ")}
}

// NewHTTPError returns an Error with the MySQL error code
func NewMySQLError(statusCode int, mysqlCode uint16, msg ...string) HTTPError {
	return httpError{Status: statusCode, MySQLCode: mysqlCode, Title: strings.Join(msg, " ")}
}

func newResponse(data interface{}, meta interface{}) *Response {
	if data == nil {
		return &Response{
			Status: http.StatusNoContent,
		}
	}

	return &Response{
		Status: http.StatusOK,
		Data:   data,
		Meta:   meta,
	}
}

// NewEmptyResponse returns an empty Response
func NewEmptyResponse() *Response {
	return &Response{}
}

type versionResponse struct {
	Version string `json:"version"`
}

// NewVersionResponse returns a Response with current version of the server
func NewVersionResponse(version string) *Response {
	return newResponse(versionResponse{version}, nil)
}

type queryMetaResponse struct {
	Headers []string `json:"headers"`
	Types   []string `json:"types"`
}

// NewQueryResponse returns a Response with table headers and row contents
func NewQueryResponse(rows []map[string]interface{}, columnNames, columnTypes []string) *Response {
	return newResponse(rows, queryMetaResponse{columnNames, columnTypes})
}
