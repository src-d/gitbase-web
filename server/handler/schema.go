package handler

import (
	"net/http"

	"github.com/src-d/gitbase-playground/server/serializer"
	"github.com/src-d/gitbase-playground/server/service"
)

// Schema returns DB schema
func Schema(db service.SQLDB) RequestProcessFunc {
	return func(r *http.Request) (res *serializer.Response, err error) {
		rows, err := db.Query("SHOW TABLES")
		if err != nil {
			return nil, err
		}

		tables := make(map[string][]serializer.Column)

		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				return nil, err
			}
			tables[table] = []serializer.Column{}
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		rows.Close()

		for table, columns := range tables {
			rows, err := db.Query("DESCRIBE TABLE " + table)
			if err != nil {
				return nil, err
			}

			for rows.Next() {
				var col serializer.Column
				if err := rows.Scan(&col.Name, &col.Type); err != nil {
					return nil, err
				}
				columns = append(columns, col)
			}

			if err := rows.Err(); err != nil {
				return nil, err
			}

			rows.Close()

			tables[table] = columns
		}

		return serializer.NewSchemaResponse(tables), nil
	}
}
