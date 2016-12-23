// Copyright 2016 Joel Scoble.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const schema = "information_schema"

type MySQLDB struct {
	DB     *sql.DB
	dbName string
}

// NewMySQLDB connects to the database's information_schema using the
// supplied username and password.  The user must have sufficient privileges.
func NewMySQLDB(server, user, password, database string) (*MySQLDB, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, schema))
	if err != nil {
		return nil, err
	}
	return &MySQLDB{
		DB:     conn,
		dbName: database,
	}, nil
}

func (m *MySQLDB) GetTables() ([]Table, error) {
	tableS := `SELECT table_schema, table_name, table_type,
	 	engine,	table_collation, table_comment
		FROM tables
		WHERE table_schema = ?`

	rows, err := m.DB.Query(tableS, m.dbName)
	if err != nil {
		return nil, err
	}
	var tables []Table
	for rows.Next() {
		var t Table
		err = rows.Scan(&t.Schema, &t.Name, &t.Typ, &t.Engine, &t.Collation, &t.Comment)
		if err != nil {
			rows.Close()
			return tables, err
		}
		tables = append(tables, t)
	}
	rows.Close()

	// go through each table and get it's columns
	columnS := `SELECT column_name, ordinal_position, column_default,
			is_nullable, data_type, character_maximum_length,
			character_octet_length, numeric_precision, numeric_scale,
			character_set_name, collation_name, column_type,
			column_key, extra, privileges,
			column_comment
		FROM columns
		WHERE table_schema = ?
			AND table_name = ?
		ORDER BY ordinal_position`

	stmt, err := m.DB.Prepare(columnS)
	if err != nil {
		return tables, err
	}
	defer stmt.Close()
	for i := range tables {
		rows, err := stmt.Query(tables[i].Schema, tables[i].Name)
		if err != nil {
			return tables, err
		}
		for rows.Next() {
			var c Column
			err = rows.Scan(&c.Name, &c.OrdinalPosition, &c.Default,
				&c.IsNullable, &c.DataType, &c.CharMaxLen,
				&c.CharOctetLen, &c.NumericPrecision, &c.NumericScale,
				&c.CharacterSet, &c.Collation, &c.Typ,
				&c.Key, &c.Extra, &c.Privileges,
				&c.Comment)
			if err != nil {
				rows.Close()
				return tables, err
			}
			tables[i].Columns = append(tables[i].Columns, c)
		}
		rows.Close()
	}
	return tables, nil
}

type Table struct {
	Name      string
	Schema    string
	Columns   []Column
	Typ       string
	Engine    string
	Collation string
	Comment   string
}

type Column struct {
	Name             string
	OrdinalPosition  uint64
	Default          sql.NullString
	IsNullable       string
	DataType         string
	CharMaxLen       sql.NullInt64
	CharOctetLen     sql.NullInt64
	NumericPrecision sql.NullInt64
	NumericScale     sql.NullInt64
	CharacterSet     sql.NullString
	Collation        sql.NullString
	Typ              string
	Key              string
	Extra            string
	Privileges       string
	Comment          string
}
