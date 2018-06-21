package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/src-d/gitbase-playground/server/service"

	"github.com/go-sql-driver/mysql"
)

type queryRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

// genericVals returns a slice of interface{}, each one a pointer to the proper
// type for each column
func genericVals(colTypes []string) []interface{} {
	columnValsPtr := make([]interface{}, len(colTypes))

	for i, colType := range colTypes {
		switch colType {
		case "BIT":
			columnValsPtr[i] = new(sql.NullBool)
		case "TIMESTAMP", "DATE", "DATETIME":
			columnValsPtr[i] = new(mysql.NullTime)
		case "INT", "MEDIUMINT", "BIGINT", "SMALLINT", "TINYINT":
			columnValsPtr[i] = new(sql.NullInt64)
		case "DOUBLE", "FLOAT":
			columnValsPtr[i] = new(sql.NullFloat64)
		case "JSON":
			columnValsPtr[i] = new([]byte)
		default: // All the text and binary variations
			columnValsPtr[i] = new(sql.NullString)
		}
	}

	return columnValsPtr
}

// Query returns a function that forwards an SQL query to gitbase and returns
// the rows as JSON
func Query(db service.SQLDB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		var queryRequest queryRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &queryRequest)
		if err != nil || queryRequest.Query == "" {
			return nil, serializer.NewHTTPError(http.StatusBadRequest,
				`Bad Request. Expected body: { "query": "SQL statement", "limit": 1234 }`)
		}

		query := addLimit(queryRequest.Query, queryRequest.Limit)
		rows, err := db.Query(query)
		if err != nil {
			return nil, dbError(err)
		}
		defer rows.Close()

		columnNames, columnTypes, err := columnsInfo(rows)
		if err != nil {
			return nil, err
		}

		columnValsPtr := genericVals(columnTypes)

		tableData := make([]map[string]interface{}, 0)

		for rows.Next() {
			if err := rows.Scan(columnValsPtr...); err != nil {
				return nil, err
			}

			colData := make(map[string]interface{}, len(columnTypes))

			for i, val := range columnValsPtr {
				colData[columnNames[i]] = nil

				switch val.(type) {
				case *sql.NullBool:
					sqlVal, _ := val.(*sql.NullBool)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.Bool
					}
				case *mysql.NullTime:
					sqlVal, _ := val.(*mysql.NullTime)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.Time
					}
				case *sql.NullInt64:
					sqlVal, _ := val.(*sql.NullInt64)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.Int64
					}
				case *sql.NullString:
					sqlVal, _ := val.(*sql.NullString)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.String
					}
				case *[]byte:
					// DatabaseTypeName JSON is used for arrays of uast nodes and
					// arrays of strings, but we don't know the exact type.
					// We try with arry of uast nodes first and any JSON later
					nodes, err := service.UnmarshallUAST(val)
					if err == nil {
						colData[columnNames[i]] = nodes
					} else {
						var data interface{}

						if err := json.Unmarshal(*val.(*[]byte), &data); err != nil {
							return nil, err
						}
						colData[columnNames[i]] = data
					}
				}
			}

			tableData = append(tableData, colData)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return serializer.NewQueryResponse(tableData, columnNames, columnTypes), nil
	}
}

// columnsInfo returns the column names and column types, or error
func columnsInfo(rows *sql.Rows) ([]string, []string, error) {
	names, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	typesStr := make([]string, len(types))
	for i, colType := range types {
		typesStr[i] = colType.DatabaseTypeName()
	}

	return names, typesStr, nil
}

var noCommentsRegexp = regexp.MustCompile(`\/\*(?s:.)*?\*\/`)
var limitRegexp = regexp.MustCompile(`\s+LIMIT\s+(\d+)$`)

// addLimit adds LIMIT to the query if it's a SELECT, avoiding '; limit'
func addLimit(query string, limit int) string {
	if limit <= 0 {
		return query
	}

	noComments := noCommentsRegexp.ReplaceAllLiteralString(query, "")

	query = strings.TrimRight(strings.TrimSpace(noComments), ";")
	upperQuery := strings.ToUpper(query)

	if strings.HasPrefix(upperQuery, "SELECT") {
		matches := limitRegexp.FindStringSubmatch(upperQuery)
		if len(matches) == 2 {
			userLimit, _ := strconv.Atoi(matches[1])
			if userLimit < limit {
				return query
			}
			query = query[:len(query)-len(matches[0])]
		}
		return fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	return query
}

// dbError transform DB error to HTTP error
func dbError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return serializer.NewMySQLError(
			http.StatusBadRequest,
			mysqlErr.Number,
			mysqlErr.Message)
	}

	return serializer.NewHTTPError(http.StatusBadRequest, err.Error())
}
