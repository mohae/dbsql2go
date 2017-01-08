// Copyright 2016-17 Joel Scoble.
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
	"bytes"
	"database/sql"
	"fmt"
	"go/format"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mohae/dbsql2go"
	"github.com/mohae/mixedcase"
)

const schema = "information_schema"

type DB struct {
	Conn    *sql.DB
	dbName  string
	tables  []dbsql2go.Tabler
	indexes []Index
	keys    []Key
	views   []dbsql2go.Viewer
}

// New connects to the database's information_schema using the supplied
// username and password.  The user must have sufficient privileges.
func New(server, user, password, database string) (dbsql2go.DBer, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, schema))
	if err != nil {
		return nil, err
	}
	return &DB{
		Conn:   conn,
		dbName: database,
	}, nil
}

func (m *DB) GetTables() error {
	tableS := `SELECT table_schema, table_name, table_type,
	 	engine,	table_collation, table_comment
		FROM tables
		WHERE table_schema = ?`

	rows, err := m.Conn.Query(tableS, m.dbName)
	if err != nil {
		return err
	}
	for rows.Next() {
		var t Table
		err = rows.Scan(&t.schema, &t.name, &t.Typ, &t.Engine, &t.collation, &t.Comment)
		if err != nil {
			rows.Close()
			return err
		}
		m.tables = append(m.tables, &t)
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

	stmt, err := m.Conn.Prepare(columnS)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i, tbl := range m.tables {
		rows, err := stmt.Query(tbl.Schema(), tbl.Name())
		if err != nil {
			return err
		}
		mTbl, ok := tbl.(*Table)
		if !ok {
			return fmt.Errorf("impossible assertion: %v is not a Table", reflect.TypeOf(tbl))
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
				return err
			}

			mTbl.Columns = append(mTbl.Columns, c)
		}
		rows.Close()
		m.tables[i] = mTbl
	}
	return nil
}

func (m *DB) Tables() []dbsql2go.Tabler {
	return m.tables
}

// GetIndexes gets the information about the databases indexes. This includes
// key column and constraint info so that indexes with constraints, i.e.
// primary keys, foreign keys, and unique can be properly identified.
//
// Any index not in the key_column_constraint is a non-unique, non-key index.
func (m *DB) GetIndexes() error {
	sel := `select TABLE_NAME, NON_UNIQUE, INDEX_SCHEMA,
		INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME,
		COLLATION, CARDINALITY, SUB_PART,
		PACKED, NULLABLE, INDEX_TYPE,
		COMMENT, INDEX_COMMENT
		from STATISTICS
		where TABLE_SCHEMA = ?
		order by TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX`

	rows, err := m.Conn.Query(sel, m.dbName)
	if err != nil {
		return err
	}
	for rows.Next() {
		var ndx Index
		err = rows.Scan(
			&ndx.TableName, &ndx.NonUnique, &ndx.IndexSchema,
			&ndx.IndexName, &ndx.SeqInIndex, &ndx.ColumnName,
			&ndx.Collation, &ndx.Cardinality, &ndx.SubPart,
			&ndx.Packed, &ndx.Nullable, &ndx.IndexType,
			&ndx.Comment, &ndx.IndexComment,
		)
		if err != nil {
			rows.Close()
			return err
		}
		m.indexes = append(m.indexes, ndx)
	}
	rows.Close()
	return nil
}

func (m *DB) GetKeys() error {
	// Get the key and constraint stuff.
	sel := `SELECT k.constraint_name, t.constraint_type, k.table_name,
	k.column_name, k.ordinal_position, k.position_in_unique_constraint,
	k.referenced_table_name, k.referenced_column_name
FROM key_column_usage AS k,
	 table_constraints AS t
WHERE k.table_schema = ?
	AND k.constraint_name = t.constraint_name
	AND k.table_name = t.table_name
GROUP BY k.table_name,
	k.constraint_name,
	k.ordinal_position`
	rows, err := m.Conn.Query(sel, m.dbName)
	if err != nil {
		return err
	}
	for rows.Next() {
		var k Key
		err = rows.Scan(
			&k.Name, &k.Type, &k.Table,
			&k.Column, &k.Seq, &k.USeq,
			&k.RefTable, &k.RefCol,
		)
		if err != nil {
			rows.Close()
			return err
		}
		m.keys = append(m.keys, k)
	}
	rows.Close()
	return nil
}

func (m *DB) GetViews() error {
	viewS := `select TABLE_NAME, VIEW_DEFINITION, CHECK_OPTION,
		IS_UPDATABLE, DEFINER, SECURITY_TYPE,
		CHARACTER_SET_CLIENT, COLLATION_CONNECTION
		from VIEWS
		where TABLE_SCHEMA = ?
		order by TABLE_NAME`

	rows, err := m.Conn.Query(viewS, m.dbName)
	if err != nil {
		return err
	}
	for rows.Next() {
		var v View
		err = rows.Scan(
			&v.TableName, &v.ViewDefinition, &v.CheckOption,
			&v.IsUpdatable, &v.Definer, &v.SecurityType,
			&v.CharacterSetClient, &v.CollationConnection,
		)
		if err != nil {
			rows.Close()
			return err
		}
		m.views = append(m.views, &v)
	}
	rows.Close()
	return nil
}

func (m *DB) Views() []dbsql2go.Viewer {
	return m.views
}

type Table struct {
	name      string
	schema    string
	Columns   []Column
	Typ       string
	Engine    sql.NullString
	collation sql.NullString
	Comment   string
	Indexes   []dbsql2go.Index
}

// Name returns the name of the table.
func (t *Table) Name() string {
	return t.name
}

// Schema returns the table's schema.
func (t *Table) Schema() string {
	return t.schema
}

// Collation returns the table's collation.
func (t *Table) Collation() string {
	if t.collation.Valid {
		return t.collation.String
	}
	return ""
}

// Go creates the struct definition an returns the resulting bytes.
// TODO: should this accept a writer instead?
func (t *Table) Go() ([]byte, error) {
	var buf bytes.Buffer
	n, err := buf.WriteString("type ")
	if err != nil {
		return nil, err
	}
	n, err = buf.WriteString(mixedcase.Exported(t.name))
	if err != nil {
		return nil, err
	}
	n, err = buf.WriteString(" struct {\n")
	if err != nil {
		return nil, err
	}

	// write the column defs
	for _, col := range t.Columns {
		err = buf.WriteByte('\t')
		if err != nil {
			return nil, err
		}
		buf.Write(col.Go())
		if err != nil {
			return nil, err
		}
		err = buf.WriteByte('\n')
		if err != nil {
			return nil, err
		}
	}
	n, err = buf.WriteString("}\n")
	if err != nil {
		return nil, err
	}
	_ = n // add short write handling
	return buf.Bytes(), nil
}

// GoFmt creates a formatted struct definition an returns the resulting bytes.
// TODO: should this accept a writer instead?
func (t *Table) GoFmt() ([]byte, error) {
	b, err := t.Go()
	if err != nil {
		return nil, err
	}
	return format.Source(b)
}

// ColumnNames returns the names of all the columns in the table.
func (t *Table) ColumnNames() []string {
	cols := make([]string, 0, len(t.Columns))
	for _, col := range t.Columns {
		cols = append(cols, col.Name)
	}
	return cols
}

// Column holds all information about the columns in a database as provided by
// MySQL's information schema.
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

func (c *Column) Go() []byte {
	n := make([]byte, 0, len(c.Name)+16) // add enough cap to handle most datatypes w/o growing
	n = append(n, []byte(mixedcase.Exported(c.Name))...)
	n = append(n, ' ')
	if c.IsNullable == "YES" {
		switch c.DataType {
		case "int", "tinyint", "smallint", "mediumint", "bigint":
			return append(n, []byte("sql.NullInt64")...)
		case "decimal":
			return append(n, []byte("sql.NullFloat64")...)
		case "timestamp", "date", "datetime":
			return append(n, []byte("mysql.NullTime")...)
		case "tinyblob", "blob", "mediumblob", "longblob",
			"tinytext", "text", "mediumtext", "longtext",
			"binary", "varbinary":
			return append(n, []byte("[]byte")...)
		case "char", "varchar", "time", "year", "enum", "set":
			return append(n, []byte("sql.NullString")...)
		default:
			return append(n, []byte(c.DataType)...)
		}
	}
	switch c.DataType {
	case "int":
		return append(n, []byte("int32")...)
	case "tinyint":
		return append(n, []byte("int8")...)
	case "smallint":
		return append(n, []byte("int16")...)
	case "mediumint":
		return append(n, []byte("int32")...)
	case "bigint":
		return append(n, []byte("int64")...)
	case "char", "varchar":
		return append(n, []byte("string")...)
	case "decimal":
		return append(n, []byte("float64")...)
	case "timestamp", "date", "datetime":
		return append(n, []byte("mysql.NullTime")...)
	case "tinyblob", "blob", "mediumblob", "longblob",
		"tinytext", "text", "mediumtext", "longtext",
		"binary", "varbinary":
		return append(n, []byte("[]byte")...)
	case "time", "year", "enum", "set":
		return append(n, []byte("string")...)
	default:
		return append(n, []byte(c.DataType)...)
	}
}

type Index struct {
	TableName    string
	NonUnique    int64
	IndexSchema  string
	IndexName    string
	SeqInIndex   int64
	ColumnName   string
	Collation    sql.NullString
	Cardinality  sql.NullInt64
	SubPart      sql.NullInt64
	Packed       sql.NullString
	Nullable     string
	IndexType    string
	Comment      sql.NullString
	IndexComment string
}

// Name returns the index's name.
func (i *Index) Name() string {
	return i.IndexName
}

// Key is data from key_column_usage and table_constraints
type Key struct {
	Name     string         // Name of the constraint
	Type     string         // Constraint type
	Table    string         // Table of the constraint
	Column   string         // Column tyhe constraint is on
	Seq      int            // Sequence number for composite constraints
	USeq     sql.NullInt64  // Position in Unique Constraint.
	RefTable sql.NullString // Table the constraint refers to for Foreign Keys
	RefCol   sql.NullString // Column on the refered to table of the constraint for Foreign Keys.
}

// ImportString returns the import string for importing the mysql db driver.
func Import() string {
	return `_ "github.com/go-sql-driver/mysql"`
}

type View struct {
	TableName           string
	ViewDefinition      string
	CheckOption         string
	IsUpdatable         string
	Definer             string
	SecurityType        string
	CharacterSetClient  string
	CollationConnection string
}

func (v *View) Name() string {
	return v.TableName
}
