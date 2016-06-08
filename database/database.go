// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package database provides functions for querying databases.
// This is an early work in progress.  This code plus the readme provides a quick start for working with postgres in golang.
package database

import (
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq" // Comment to justify blank import
)

// GetWidgets provides ...
func GetWidgets() ([]byte, error) {

	// Get keys from DB
	db, err := sql.Open("postgres", "dbname=dev user=dev password=dev host=localhost sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	data, err := doSelect(db)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil

}

// Maybe change this to take a selectStatement as a variable
// It takes parameters blah blah and returns blah blah.
func doSelect(db *sql.DB) ([]map[string]interface{}, error) {
	stmt, err := db.Prepare(`select * from widgets`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)

	var tableData []map[string]interface{}
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil

}
