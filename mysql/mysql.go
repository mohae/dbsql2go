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
	"io"
	"reflect"
	"unicode"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mohae/dbsql2go"
	"github.com/mohae/int2word"
	"github.com/mohae/mixedcase"
)

const (
	schema                 = "information_schema"
	viewType               = "VIEW"
	selectPKComment        = "Select SELECTs the row from %s that corresponds with the struct's primary key and populates the struct with the SELECTed data. Any error that occurs will be returned."
	selectPKInRangeComment = "%sSelectInRange%s SELECTs a range of rows from the %s table whose PK values are within the specified range and returns a slice of %s structs. The range values are %s. %s args must be passed for the values of the query's range boundaries in the WHERE clause. The WHERE clause is in the form of %q. If there is an error, the error will be returned and the results slice will be nil."
	deletePKComment        = "Delete DELETEs the row from %s that corresponds with the struct's primary key, if there is any. The number of rows DELETEd is returned. If an error occurs during the DELETE, an error will be returned along with 0."
	insertPKComment        = "Insert INSERTs the data in the struct into %s. The ID from the INSERT, if applicable, is returned. If an error occurs that is returned along with a 0."
	updatePKComment        = "Update UPDATEs the row in %s that corresponds with the struct's key values. The number of rows affected by the update will be returned. If an error occurs, the error will be returned along with 0."
)

type DB struct {
	Conn        *sql.DB
	Name        string
	tables      []dbsql2go.Tabler
	indexes     []Index
	constraints []Constraint
	views       []dbsql2go.Viewer
}

// New connects to the database's information_schema using the supplied
// username and password.  The user must have sufficient privileges.
func New(server, user, password, database string) (dbsql2go.DBer, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, schema))
	if err != nil {
		return nil, err
	}
	return &DB{
		Conn: conn,
		Name: database,
	}, nil
}

// Get retrieves all of the table, view, index, and constraint info for a
// database. The tables will have information about their constraints and
// indexes. None of the other Get or Update methods need to be called when
// using this method.
func (m *DB) Get() error {
	err := m.GetTables()
	if err != nil {
		return err
	}

	err = m.GetViews()
	if err != nil {
		return err
	}

	err = m.GetIndexes()
	if err != nil {
		return err
	}

	err = m.GetConstraints()
	if err != nil {
		return err
	}

	m.UpdateTableIndexes()
	err = m.UpdateTableConstraints()
	if err != nil {
		return err
	}
	return nil
}

func (m *DB) GetTables() error {
	tableS := `SELECT table_schema, table_name, table_type,
	 	engine,	table_collation, table_comment
		FROM tables
		WHERE table_schema = ?`

	rows, err := m.Conn.Query(tableS, m.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		t := NewTable()
		err = rows.Scan(&t.schema, &t.name, &t.Typ, &t.Engine, &t.collation, &t.Comment)
		if err != nil {
			rows.Close()
			return err
		}
		m.tables = append(m.tables, t)
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
			// set the column's corresponding Go field name
			c.SetFieldName()
			mTbl.columns = append(mTbl.columns, c)
		}
		rows.Close()
		mTbl.sqlInf.Table = tbl.Name()
		mTbl.structName = mixedcase.Exported(tbl.Name())
		r, _ := utf8.DecodeRuneInString(tbl.StructName())
		mTbl.r = unicode.ToLower(r)
		m.tables[i] = mTbl
	}
	return nil
}

// Tables returns information about all of the tables in a databasse; this
// includes views but not view specific information like its definition.
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

	rows, err := m.Conn.Query(sel, m.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		var ndx Index
		err = rows.Scan(
			&ndx.Table, &ndx.NonUnique, &ndx.Schema,
			&ndx.name, &ndx.SeqInIndex, &ndx.Column,
			&ndx.Collation, &ndx.Cardinality, &ndx.SubPart,
			&ndx.Packed, &ndx.Nullable, &ndx.Type,
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

func (m *DB) GetConstraints() error {
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
	rows, err := m.Conn.Query(sel, m.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		var c Constraint
		err = rows.Scan(
			&c.Name, &c.Type, &c.Table,
			&c.Column, &c.Seq, &c.USeq,
			&c.RefTable, &c.RefCol,
		)
		if err != nil {
			rows.Close()
			return err
		}
		m.constraints = append(m.constraints, c)
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

	rows, err := m.Conn.Query(viewS, m.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		var v View
		err = rows.Scan(
			&v.Table, &v.ViewDefinition, &v.CheckOption,
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

// UpdateTableConstraints updates the Tables with their respective Constraint
// information. The Constraints must be retrieved first or nothing will be
// done.
func (m *DB) UpdateTableConstraints() error {
	// Map the retrieved constraints back to their respective tables. There may
	// be multiple constraints per table and multiple rows per constraint.
	var prior Constraint
	var c dbsql2go.Constraint
	for i, v := range m.constraints {
		if v.Table == prior.Table && v.Name == prior.Name { // if this is just another row for the same constraint, add the info
			c.Columns = append(c.Columns, v.Column)
			c.Fields = append(c.Fields, fieldName(v.Column))
			if v.RefCol.Valid {
				c.RefColumns = append(c.RefColumns, v.RefCol.String)
				c.RefFields = append(c.RefFields, fieldName(v.RefCol.String))
			}
			prior = v
			continue
		}
		// if this is the first entry; don't add the index
		if i == 0 {
			goto process
		}
		// find the table and add this index to it
		for j := 0; j < len(m.tables); j++ {
			if m.tables[j].Name() != c.Table {
				continue
			}
			m.tables[j].(*Table).constraints = append(m.tables[j].(*Table).constraints, c)
			if c.Type == dbsql2go.PK { // if the constraint type is pk, set the index for pk
				m.tables[j].(*Table).pk = len(m.tables[j].(*Table).constraints) - 1
			}
			break
		}
	process:
		typ, err := dbsql2go.ParseConstraintType(v.Type)
		if err != nil {
			return err
		}
		c = dbsql2go.Constraint{Type: typ, Name: v.Name, Table: v.Table, Columns: []string{v.Column}, Fields: []string{fieldName(v.Column)}}
		if v.RefTable.Valid {
			c.RefTable = v.RefTable.String
		}
		if v.RefCol.Valid {
			c.RefColumns = append(c.RefColumns, v.RefCol.String)
			c.RefFields = append(c.RefFields, fieldName(v.RefCol.String))
		}
		prior = v
	}
	// handle the final element
	if prior.Name != "" {
		// find the table and add this index to it
		for i := 0; i < len(m.tables); i++ {
			if m.tables[i].Name() != prior.Table {
				continue
			}
			m.tables[i].(*Table).constraints = append(m.tables[i].(*Table).constraints, c)
			break
		}
	}
	return nil
}

// UpdateTableIndexes updates the Tables with their respective Index information.
// The Indexes must be retrieved first or nothing will be done.
func (m *DB) UpdateTableIndexes() {
	// Map the retrieved indexes back to their respective tables. There may be
	// multiple indexes per table and multiple rows per index.
	var prior Index
	var ndx dbsql2go.Index
	for i, v := range m.indexes {
		if v.Table == prior.Table && v.name == prior.name { // if this is just another row for the same index, add the info
			ndx.Columns = append(ndx.Columns, v.Column)
			prior = v
			continue
		}
		// if this is the first entry; don't add the index
		if i == 0 {
			goto process
		}
		// find the table and add this index to it
		for j := 0; j < len(m.tables); j++ {
			if m.tables[j].Name() != prior.Table {
				continue
			}
			m.tables[j].(*Table).indexes = append(m.tables[j].(*Table).indexes, ndx)
			break
		}
	process:
		ndx = dbsql2go.Index{Type: v.Type, Name: v.name, Table: v.Table, Columns: []string{v.Column}}
		if v.name == "PRIMARY" {
			ndx.Primary = true
		}
		prior = v
	}
	// handle the final element
	if prior.name != "" {
		// find the table and add this index to it
		for i := 0; i < len(m.tables); i++ {
			if m.tables[i].Name() != prior.Table {
				continue
			}
			m.tables[i].(*Table).indexes = append(m.tables[i].(*Table).indexes, ndx)
			break
		}
	}
}

type Table struct {
	//db          *DB
	name        string
	r           rune   // the first letter of the name, in lower-case. Used as the receiver name.
	structName  string // the name of the struct for this table
	schema      string
	columns     []Column
	Typ         string
	Engine      sql.NullString
	collation   sql.NullString
	Comment     string
	indexes     []dbsql2go.Index
	constraints []dbsql2go.Constraint
	pk          int               // index of the pk constraint in constraints, if there is one
	sqlInf      dbsql2go.TableSQL // caches all columns for the table for SQL generation
	buf         bytes.Buffer      // buffer for holding generated stuff; this is not thread-safe
}

// NewTable creates a new Table. It is intended to ensure that the Table is
// ready for usage.
func NewTable() *Table {
	// set pk to a negative value to indicate there isn't one.
	return &Table{pk: -1}
}

// Name returns the name of the table.
func (t *Table) Name() string {
	return t.name
}

// StructName returns the name of the Go struct for this table.
func (t *Table) StructName() string {
	return t.structName
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

// Definition creates the struct definition and appends it to the  internal
// buffer. The number of bytes written to the buffer is returned along with
// an error, if one occurs.
// TODO: should this accept a writer instead
func (t *Table) Definition(w io.Writer) error {
	_, err := w.Write([]byte("type "))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(t.structName))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(" struct {\n"))
	if err != nil {
		return err
	}

	// write the column defs
	for _, col := range t.columns {
		_, err = w.Write([]byte{'\t'})
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(col.Go()))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte("}\n"))
	if err != nil {
		return err
	}
	return nil
}

// Go creats the struct definition and methods for handling single row
// SQL queries that the struct will use. A struct represents one row of data.
// Any operations that result in more than one row are handled by something
// other than the table's struct.
func (t *Table) Go(w io.Writer) error {
	// generate the struct def
	err := t.Definition(w)
	if err != nil {
		return err
	}

	// add the select method
	err = t.SelectPKMethod(w)
	if err != nil {
		return err
	}

	// add the delete method
	err = t.DeletePKMethod(w)
	if err != nil {
		return err
	}

	// add the insert method
	err = t.InsertMethod(w)
	if err != nil {
		return err
	}

	err = t.UpdateMethod(w)
	if err != nil {
		return err
	}

	_, err = t.SelectInRangeFunc(w)
	return nil
}

// GoFmt creates a formatted struct definition and methods and returns the
// resulting bytes.
func (t *Table) GoFmt(w io.Writer) error {
	// use a buffer for the defintion so that it can be formatted before writing
	var buf bytes.Buffer
	err := t.Go(&buf)
	if err != nil {
		return fmt.Errorf("%s: create definition: %s", t.name, err)
	}

	// format the definition
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("%s: format definition: %s", t.name, err)
	}
	// write the definition
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("%s: write definition: %s", t.name, err)
	}
	return nil
}

// ColumnNames returns the names of all the columns in the table.
func (t *Table) ColumnNames() []string {
	Columns := make([]string, 0, len(t.columns))
	for _, col := range t.columns {
		Columns = append(Columns, col.Name)
	}
	return Columns
}

// NonPKColumnNames returns the names of all the non-pk columns in the table
// TODO: is this still necessary?
func (t *Table) NonPKColumnNames() []string {
	pk := t.PK()
	if pk == nil {
		return t.ColumnNames() // if there isn't a pk on this table, return all columns
	}

	Columns := make([]string, 0, len(t.columns))
	for _, col := range t.columns {
		var pkCol bool
		// this isn't optimal but good enough considering number of PK Columns will be low, if any
		for _, v := range pk.Columns {
			if v == col.Name {
				pkCol = true
				break
			}
		}
		if !pkCol {
			Columns = append(Columns, col.Name)
		}
	}
	return Columns
}

// NonAutoIncrementColumnNames returns the names of all the auto-increment
// columns in the table
// TODO: is this still necessary?
func (t *Table) NonAutoIncrementColumnNames() []string {
	Columns := make([]string, 0, len(t.columns))
	for _, col := range t.columns {
		if col.Extra == "auto_increment" {
			continue
		}
		Columns = append(Columns, col.Name)
	}
	return Columns
}

// Indexes returns information on all of the tables indexes.
func (t *Table) Indexes() []dbsql2go.Index {
	return t.indexes
}

// Constraints returns information on all of the tables keys/constraints.
func (t *Table) Constraints() []dbsql2go.Constraint {
	return t.constraints
}

// IsView returns whether or not this table is actually a view.
func (t *Table) IsView() bool {
	if t.Typ == viewType {
		return true
	}
	return false
}

// SelectPKMethod generates the method for selecting a table row using its PK.
func (t *Table) SelectPKMethod(w io.Writer) error {
	pk := t.PK()
	if pk == nil {
		// nothing to do
		return nil
	}

	_, err := w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// write the comment
	c, err := dbsql2go.StringToComments(fmt.Sprintf(selectPKComment, t.name), 80)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(c))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fmt.Sprintf("func(%c *%s) Select(db *sql.DB) error {\n\terr := db.QueryRow(\"", t.r, t.structName)))
	if err != nil {
		return err
	}

	err = t.SelectSQLPK(w)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\", "))
	if err != nil {
		return err
	}

	for i, v := range pk.Fields {
		if i == 0 {
			_, err = w.Write([]byte(fmt.Sprintf("%c.%s", t.r, v)))
			if err != nil {
				return err
			}
			continue
		}
		_, err = w.Write([]byte(fmt.Sprintf(", %c.%s", t.r, v)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(").Scan("))
	if err != nil {
		return err
	}

	// buld the struct field stuff
	for i, v := range t.columns {
		if i == 0 {
			_, err = w.Write([]byte(fmt.Sprintf("&%c.%s", t.r, v.fieldName)))
			if err != nil {
				return err
			}
			continue
		}
		_, err = w.Write([]byte(fmt.Sprintf(", &%c.%s", t.r, v.fieldName)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(")\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}"))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	return nil
}

// SelectSQLPK returns a SELECT statement for the table that selects all the
// table columns using the tables PK. If the table does not have a PK, a nil
// will be returned and the error will also be nil as this is not an error
// state.
func (t *Table) SelectSQLPK(w io.Writer) error {
	pk := t.PK()
	if pk == nil { // the table doesn't have a primary key; this is not an error.
		return nil
	}

	// set up the relevant infor for the SQL generation; Table is already set.
	t.sqlInf.Columns = t.ColumnNames()
	t.sqlInf.WhereColumns = pk.Columns
	err := dbsql2go.SelectSQL.Execute(w, t.sqlInf)
	if err != nil {
		return err
	}
	return nil
}

// SelectRangeSQL creates in range SELECT funcs for the table if it has a
// primary key. Tables without priamry keys wiill have nothing written to the
// writer and 0 will be returned for the number of bytes written along with nil
// for the error. Any error encountered is written along with the number of
// bytes for the table.
func (t *Table) SelectInRangeFunc(w io.Writer) (n int64, err error) {
	// If the table has a primary key
	pk := t.PK()
	if pk == nil { // If no primary key return 0 for bytes written and nil for the error.
		return 0, nil
	}

	// Prepare the Table Information for the SQL
	t.sqlInf.Columns = t.ColumnNames()
	// Reset the where info
	t.sqlInf.WhereColumns = t.sqlInf.WhereColumns[:0]
	t.sqlInf.WhereComparisonOps = t.sqlInf.WhereComparisonOps[:0]
	t.sqlInf.WhereConditions = t.sqlInf.WhereConditions[:0]

	// for Where columns, each pk column is used twice to set up >= <=.
	for i, col := range pk.Columns {
		if i != 0 {
			t.sqlInf.WhereConditions = append(t.sqlInf.WhereConditions, "AND")
		}
		t.sqlInf.WhereColumns = append(t.sqlInf.WhereColumns, col)
		t.sqlInf.WhereColumns = append(t.sqlInf.WhereColumns, col)
		t.sqlInf.WhereConditions = append(t.sqlInf.WhereConditions, "AND")
	}

	n, err = t.selectInRangeInclusive(w)
	if err != nil {
		return 0, nil
	}

	return n, nil
}

func (t *Table) selectInRangeInclusive(w io.Writer) (n int64, err error) {

	// Set the Comparison ops for inclusive
	for i := 0; i < len(t.sqlInf.WhereColumns)/2; i++ {
		t.sqlInf.WhereComparisonOps = append(t.sqlInf.WhereComparisonOps, ">=")
		t.sqlInf.WhereComparisonOps = append(t.sqlInf.WhereComparisonOps, "<=")
	}
	// reset the buffer: everything gets written to the buffer first
	t.buf.Reset()

	err = dbsql2go.SelectAndOrWhereComment.Execute(&t.buf, t.sqlInf)
	if err != nil {
		return 0, err
	}

	num := int2word.Capitalized(int64(len(t.sqlInf.WhereColumns)))
	c, err := dbsql2go.StringToComments(fmt.Sprintf(selectPKInRangeComment, t.structName, dbsql2go.TitleInclusive, t.name, t.structName, dbsql2go.LowerInclusive, num, t.buf.String()), 80)
	if err != nil {
		return 0, err
	}

	t.buf.Reset()

	err = t.buf.WriteByte('\n')
	if err != nil {
		return 0, err
	}

	_, err = t.buf.WriteString(c)
	if err != nil {
		return 0, err
	}

	_, err = t.buf.WriteString(fmt.Sprintf("func %sSelectInRangeInclusive(db *sql.DB, args ...interface{}) (results []%s, err error) {\n\trows, err := db.Query(\"", t.structName, t.structName))
	if err != nil {
		return 0, err
	}

	// write the sql
	err = dbsql2go.SelectAndOrSQL.Execute(&t.buf, t.sqlInf)
	if err != nil {
		return 0, err
	}

	_, err = t.buf.WriteString("\", args...)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tdefer rows.Close()\n\n\tfor rows.Next() {\n\t\tvar ")
	if err != nil {
		return 0, err
	}

	_, err = t.buf.WriteString(fmt.Sprintf("%c %s\n\t\terr = rows.Scan(", t.r, t.structName))
	if err != nil {
		return 0, err
	}

	_, err = t.buf.WriteString(fmt.Sprintf("&%c.%s", t.r, t.columns[0].fieldName))
	if err != nil {
		return 0, err
	}

	for i := 1; i < len(t.columns); i++ {
		_, err = t.buf.Write([]byte(fmt.Sprintf(", &%c.%s", t.r, t.columns[i].fieldName)))
		if err != nil {
			return 0, err
		}
	}

	_, err = t.buf.WriteString(")\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tresults = append(results, ")
	if err != nil {
		return 0, err
	}

	_, err = t.buf.WriteString(fmt.Sprintf("%c)\n\t}\n\n\treturn results, nil\n}\n", t.r))
	if err != nil {
		return 0, err
	}

	return t.buf.WriteTo(w)
}

// DeletePKMethod generates the method for deleting a table row using its PK.
func (t *Table) DeletePKMethod(w io.Writer) error {
	pk := t.PK()
	if pk == nil {
		// nothing to do
		return nil
	}

	_, err := w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// write the comment
	c, err := dbsql2go.StringToComments(fmt.Sprintf(deletePKComment, t.name), 80)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(c))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fmt.Sprintf("func(%c *%s) Delete(db *sql.DB) (n int64, err error) {\n\tres, err := db.Exec(\"", t.r, t.structName)))
	if err != nil {
		return err
	}

	err = t.DeleteSQLPK(w)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\", "))
	if err != nil {
		return err
	}

	for i, v := range pk.Fields {
		if i == 0 {
			_, err = w.Write([]byte(fmt.Sprintf("%c.%s", t.r, v)))
			if err != nil {
				return err
			}
			continue
		}
		_, err = w.Write([]byte(fmt.Sprintf(", %c.%s", t.r, v)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(")\n\tif err != nil {\n\t\treturn 0, err\n\t}\n\treturn res.RowsAffected()\n}\n"))
	if err != nil {
		return err
	}

	return nil
}

// DeleteSQLPK returns a DELETE statement for the table that deletes a row
// using the tables PK. If the table does not have a PK, no SQL will be
// generated and a nil will be returned as this is not an error state.
func (t *Table) DeleteSQLPK(w io.Writer) error {
	pk := t.PK()
	if pk == nil { // the table doesn't have a primary key; this is not an error.
		return nil
	}
	t.sqlInf.WhereColumns = pk.Columns
	err := dbsql2go.DeleteSQL.Execute(w, t.sqlInf)
	if err != nil {
		return err
	}
	return nil
}

// InsertMethod generates the method for inserting the Table's data into the
// db table.
func (t *Table) InsertMethod(w io.Writer) error {
	pk := t.PK()
	if pk == nil {
		// nothing to do
		return nil
	}

	_, err := w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// write the comment
	c, err := dbsql2go.StringToComments(fmt.Sprintf(insertPKComment, t.name), 80)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(c))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fmt.Sprintf("func(%c *%s) Insert(db *sql.DB) (id int64, err error) {\n\tres, err := db.Exec(\"", t.r, t.structName)))
	if err != nil {
		return err
	}

	err = t.InsertSQL(w)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\", "))
	if err != nil {
		return err
	}

	// buld the struct field stuff: skip the pk Columns and only use the fields
	// that have corresponding columns in sqlInf
	var j int // index into the sqlInf Columns
	for _, v := range t.columns {
		// skip if this column isn't in sqlInf; columns are in same order
		if v.Name != t.sqlInf.Columns[j] {
			continue
		}
		j++         // point to next column
		if j == 1 { // if this is the first element added, don't prefix with ', '
			_, err = w.Write([]byte(fmt.Sprintf("&%c.%s", t.r, v.fieldName)))
			if err != nil {
				return err
			}
			continue
		}
		_, err = w.Write([]byte(fmt.Sprintf(", &%c.%s", t.r, v.fieldName)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(")\n\tif err != nil {\n\t\treturn 0, err\n\t}\n\treturn res.LastInsertID()\n}\n"))
	if err != nil {
		return err
	}

	return nil
}

// InsertSQL returns an INSERT statement for the table.
func (t *Table) InsertSQL(w io.Writer) error {
	// don't generate sql for views
	if t.IsView() {
		return nil
	}

	// set up the relevant infor for the SQL generation; Table is already set.
	t.sqlInf.Columns = t.NonAutoIncrementColumnNames()

	err := dbsql2go.InsertSQL.Execute(w, t.sqlInf)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMethod generates the method for updatating a table row using its PK.
func (t *Table) UpdateMethod(w io.Writer) error {
	pk := t.PK()
	if pk == nil {
		// nothing to do
		return nil
	}

	_, err := w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// write the comment
	c, err := dbsql2go.StringToComments(fmt.Sprintf(updatePKComment, t.name), 80)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(c))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(fmt.Sprintf("func(%c *%s) Update(db *sql.DB) (n int64, err error) {\n\tres, err := db.Exec(\"", t.r, t.structName)))
	if err != nil {
		return err
	}

	err = t.UpdateSQL(w)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\", "))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(""))
	if err != nil {
		return err
	}

	var j int //keep track of what's been added

	// buld the struct field stuff
	for _, v := range t.columns {
		var b bool
		// if the column isn't in the list of columns added to the query, skip it
		for _, col := range t.sqlInf.Columns {
			if v.Name == col {
				b = true
				break
			}
		}
		if !b {
			continue
		}
		j++
		if j == 1 {
			_, err = w.Write([]byte(fmt.Sprintf("&%c.%s", t.r, v.fieldName)))
			if err != nil {
				return err
			}
			continue
		}
		_, err = w.Write([]byte(fmt.Sprintf(", &%c.%s", t.r, v.fieldName)))
		if err != nil {
			return err
		}
	}

	for _, v := range pk.Fields {
		_, err = w.Write([]byte(fmt.Sprintf(", &%c.%s", t.r, v)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte(")\n\tif err != nil {\n\t\treturn 0, err\n\t}\n\treturn res.RowsAffected()\n}"))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	return nil
}

// UpdateSQLPK returns an UPDATE statement for the table that updates all
// non-auto increment columns using the table's pk in the WHERE clause. If
// the table does not have a PK, no UPDATE statement will be generated and a
// nil will be returned as this is not an error state.
func (t *Table) UpdateSQL(w io.Writer) error {
	pk := t.PK()
	if pk == nil { // the table doesn't have a primary key; this is not an error.
		return nil
	}

	// set up the relevant infor for the SQL generation; Table is already set.
	t.sqlInf.Columns = t.NonAutoIncrementColumnNames()
	t.sqlInf.WhereColumns = pk.Columns
	err := dbsql2go.UpdateSQL.Execute(w, t.sqlInf)
	if err != nil {
		return err
	}
	return nil
}

// PK returns a tables primary key information, if it has a primary key, or
// nil if it doesn't have a primary key
func (t *Table) PK() *dbsql2go.Constraint {
	if t.pk < 0 { // if the index is negative this table doesn't have a pk
		return nil
	}
	return &t.constraints[t.pk]
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
	fieldName        string
}

func (c *Column) Go() []byte {
	n := make([]byte, 0, len(c.Name)+16) // add enough cap to handle most datatypes w/o growing
	n = append(n, []byte(c.fieldName)...)
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

// SetFieldName sets the column's field name; the name of the field in the
// table struct in which this column's value will be put.
func (c *Column) SetFieldName() {
	c.fieldName = fieldName(c.Name)
}

type Index struct {
	name         string
	Type         string
	Schema       string
	NonUnique    int64
	SeqInIndex   int64
	Table        string
	Column       string
	Collation    sql.NullString
	Cardinality  sql.NullInt64
	SubPart      sql.NullInt64
	Packed       sql.NullString
	Nullable     string
	Comment      sql.NullString
	IndexComment string
}

// Name returns the index's name.
func (i *Index) Name() string {
	return i.name
}

// Constraint is data from key_column_usage and table_constraints
type Constraint struct {
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
	Table               string
	ViewDefinition      string
	CheckOption         string
	IsUpdatable         string
	Definer             string
	SecurityType        string
	CharacterSetClient  string
	CollationConnection string
}

func (v *View) Name() string {
	return v.Table
}

// fieldName makes returns an exported Go fieldName from the received string.
func fieldName(s string) string {
	return mixedcase.Exported(s)
}
