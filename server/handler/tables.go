package handler

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/src-d/gitbase-playground/server/serializer"
)

// Tables returns a function that calls /query with the SQL "SHOW TABLES"
func Tables(db *sql.DB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		req, _ := http.NewRequest("POST", "/query",
			strings.NewReader(`{ "query": "SHOW TABLES" }`))

		return Query(db)(req)
	}
}
