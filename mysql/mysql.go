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
	"github.com/mohae/dbsql2go"
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
	s := `SELECT table_name, table_type, engine,
	 	table_collation, table_comment
		FROM tables
		WHERE table_schema = ?`

	rows, err := m.DB.Query(s, m.dbName)
	if err != nil {
		return nil, err
	}
	var tables []Table
	defer rows.Close()
	for rows.Next() {
		var t Table
		err = rows.Scan(&t.name, &t.Type, &t.Engine, &t.Collation, &t.Comment)
		if err != nil {
			return tables, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

type Table struct {
	name      string
	columns   []dbsql2go.Column
	Type      string
	Engine    string
	Collation string
	Comment   string
}
