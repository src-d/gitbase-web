package handler

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"

	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/src-d/gitbase-playground/server/service"
)

// Export returns a function that forwards an SQL query to gitbase and returns
// the rows as CSV file
func Export(db service.SQLDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func(w http.ResponseWriter, r *http.Request) error {
			query := r.URL.Query().Get("query")
			if query == "" {
				return serializer.NewHTTPError(http.StatusBadRequest,
					`Bad Request. Query can't be empty.`)
			}

			rows, err := db.Query(query)
			if err != nil {
				return dbError(err)
			}
			defer rows.Close()

			columnNames, columnTypes, err := columnsInfo(rows)
			if err != nil {
				return err
			}

			columnValsPtr := genericVals(columnTypes)

			w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
			w.Header().Set("Content-Type", "text/csv")
			csvWriter := csv.NewWriter(w)

			csvWriter.Write(columnNames)

			for rows.Next() {
				if err := rows.Scan(columnValsPtr...); err != nil {
					return err
				}

				record := make([]string, len(columnTypes))

				for i, val := range columnValsPtr {
					switch v := val.(type) {
					case *sql.NullBool:
						sqlVal, _ := val.(*sql.NullBool)
						if sqlVal.Valid {
							record[i] = strconv.FormatBool(sqlVal.Bool)
						}
					case *mysql.NullTime:
						sqlVal, _ := val.(*mysql.NullTime)
						if sqlVal.Valid {
							b, err := sqlVal.Time.MarshalText()
							if err != nil {
								return err
							}
							record[i] = string(b)
						}
					case *sql.NullInt64:
						sqlVal, _ := val.(*sql.NullInt64)
						if sqlVal.Valid {
							record[i] = strconv.FormatInt(sqlVal.Int64, 10)
						}
					case *sql.NullString:
						sqlVal, _ := val.(*sql.NullString)
						if sqlVal.Valid {
							record[i] = sqlVal.String
						}
					case *[]byte:
						// DatabaseTypeName JSON is used for arrays of uast nodes and
						// arrays of strings, but we don't know the exact type.
						// We try with arry of uast nodes first and any JSON later
						nodes, err := service.UnmarshallUAST(val)
						if err == nil {
							b, err := json.Marshal(nodes)
							if err != nil {
								return err
							}
							record[i] = string(b)
							continue
						}

						record[i] = string(*v)
					}
				}

				if err := csvWriter.Write(record); err != nil {
					return err
				}
			}

			if err := rows.Err(); err != nil {
				return err
			}

			if err := csvWriter.Error(); err != nil {
				return err
			}

			csvWriter.Flush()

			return nil
		}(w, r)

		if err != nil {
			if httpError, ok := err.(serializer.HTTPError); ok {
				http.Error(w, httpError.Error(), httpError.StatusCode())
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
