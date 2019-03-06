package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/src-d/gitbase-web/server/serializer"
	"github.com/src-d/gitbase-web/server/service"

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
		default:
			// All the text and binary variations. For some reason BLOB is actually
			// returned as TEXT
			columnValsPtr[i] = new(sql.NullString)
		}
	}

	return columnValsPtr
}

// Query returns a function that forwards an SQL query to gitbase and returns
// the rows as JSON
func Query(db service.SQLDB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		var queryReq queryRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &queryReq)
		if err != nil || queryReq.Query == "" {
			return nil, serializer.NewHTTPError(http.StatusBadRequest,
				`Bad Request. Expected body: { "query": "SQL statement", "limit": 1234 }`)
		}

		// go-sql-driver/mysql QueryContext stops waiting for the query results on
		// context cancel, but it does not actually cancel the query on the server

		c := make(chan error, 1)

		conn, err := db.Conn(r.Context())
		if err != nil {
			return nil, fmt.Errorf("failed to get a DB connection: %s", err)
		}
		defer conn.Close()

		connID, err := getConnID(conn)
		if err != nil {
			return nil, fmt.Errorf("failed to get connection id: %s", err)
		}

		var resp *serializer.Response
		go func() {
			resp, err = queryContext(r.Context(), conn, queryReq)
			c <- err
		}()

		// It may happen that the QueryContext returns with an error because of
		// context cancellation. In this case, the select may enter on the second
		// case. We check if the context was cancelled with Err() instead of Done()
		select {
		case <-r.Context().Done():
		case err = <-c:
		}

		if r.Context().Err() != nil {
			db.Exec(fmt.Sprintf("KILL %d", connID))
			return nil, dbError(r.Context().Err())
		}

		if err != nil {
			return nil, dbError(err)
		}

		return resp, nil
	}
}

func queryContext(ctx context.Context, conn *sql.Conn, queryReq queryRequest) (*serializer.Response, error) {
	query, limitSet := addLimit(queryReq.Query, queryReq.Limit)

	var rows *sql.Rows

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
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

		colData, err := columnsData(columnNames, columnTypes, columnValsPtr)
		if err != nil {
			return nil, err
		}

		tableData = append(tableData, colData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return serializer.NewQueryResponse(
		tableData, columnNames, columnTypes, limitSet, queryReq.Limit), nil
}

func getConnID(conn *sql.Conn) (uint32, error) {
	const connIDQuery = "SELECT CONNECTION_ID()"
	var connID uint32

	row := conn.QueryRowContext(context.Background(), connIDQuery)
	if row == nil {
		return 0, fmt.Errorf("failed to execute %q", connIDQuery)
	}

	if err := row.Scan(&connID); err != nil {
		return 0, fmt.Errorf("failed to scan the results of %q: %s", connIDQuery, err)
	}

	return connID, nil
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

func columnsData(
	columnNames []string,
	columnTypes []string,
	columnValsPtr []interface{},
) (map[string]interface{}, error) {
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
		case *sql.NullFloat64:
			sqlVal, _ := val.(*sql.NullFloat64)
			if sqlVal.Valid {
				colData[columnNames[i]] = sqlVal.Float64
			}
		case *sql.NullString:
			// DatabaseTypeName TEXT is used for text or blobs. We try
			// to parse as UAST first
			sqlVal, _ := val.(*sql.NullString)
			if sqlVal.Valid {
				nodes, err := service.UnmarshalNodes([]byte(sqlVal.String))
				if err == nil && nodes != nil {
					colData[columnNames[i]] = nodes
					colData["__"+columnNames[i]+"-protobufs"] = []byte(sqlVal.String)
				} else {
					colData[columnNames[i]] = sqlVal.String
				}
			}
		case *[]byte:
			// DatabaseTypeName JSON is used for arrays of strings
			var data interface{}

			if err := json.Unmarshal(*val.(*[]byte), &data); err != nil {
				return nil, err
			}
			colData[columnNames[i]] = data
		}
	}

	return colData, nil
}

var noCommentsRegexp = regexp.MustCompile(`\/\*(?s:.)*?\*\/`)
var limitRegexp = regexp.MustCompile(`\s+LIMIT\s+(\d+)$`)

// addLimit adds LIMIT to the query if it's a SELECT, avoiding '; limit'
// returns true if the limit was applied
func addLimit(query string, limit int) (string, bool) {
	if limit <= 0 {
		return query, false
	}

	noComments := noCommentsRegexp.ReplaceAllLiteralString(query, "")

	query = strings.TrimSpace(strings.TrimRight(strings.TrimSpace(noComments), ";"))
	upperQuery := strings.ToUpper(query)

	if strings.HasPrefix(upperQuery, "SELECT") {
		matches := limitRegexp.FindStringSubmatch(upperQuery)
		if len(matches) == 2 {
			userLimit, _ := strconv.Atoi(matches[1])
			if userLimit <= limit {
				return query, false
			}
			query = query[:len(query)-len(matches[0])]
		}
		return fmt.Sprintf("%s LIMIT %d", query, limit), true
	}

	return query, false
}

// dbError transform DB error to HTTP error
func dbError(err error) error {
	if err == context.Canceled {
		return err
	}

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return serializer.NewMySQLError(
			http.StatusBadRequest,
			mysqlErr.Number,
			mysqlErr.Message)
	}

	return serializer.NewHTTPError(http.StatusBadRequest, err.Error())
}
