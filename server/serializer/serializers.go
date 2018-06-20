package serializer

import (
	"net/http"
	"strings"

	"github.com/src-d/gitbase-playground/server/service"
	enry "gopkg.in/src-d/enry.v1"
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

// NewMySQLError returns an Error with the MySQL error code
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

// Column describes a table column in DB
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Table struct describes a schema of one table in DB
type Table struct {
	Table   string   `json:"table"`
	Columns []Column `json:"columns"`
}

// NewSchemaResponse returns a Response with tables schema
func NewSchemaResponse(tables map[string][]Column) *Response {
	var res []Table
	for table, columns := range tables {
		res = append(res, Table{
			Table:   table,
			Columns: columns,
		})
	}
	return newResponse(res, nil)
}

// NewParseResponse returns a Response with UAST
func NewParseResponse(resp *service.ParseResponse) *Response {
	return newResponse(resp, nil)
}

// NewDetectLangResponse returns a Response with detected language
func NewDetectLangResponse(lang string, langType enry.Type) *Response {
	return newResponse(struct {
		Language string `json:"language"`
		Type     int    `json:"type"`
	}{lang, int(langType)}, nil)
}
