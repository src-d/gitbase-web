package handler

import (
	"net/http"
	"strings"

	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/src-d/gitbase-playground/server/service"
)

// Tables returns a function that calls /query with the SQL "SHOW TABLES"
func Tables(db service.SQLDB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		req, _ := http.NewRequest("POST", "/query",
			strings.NewReader(`{ "query": "SHOW TABLES" }`))

		return Query(db)(req)
	}
}
