package dbsql2go

// TableSQL is used to describe basic components of a sql statement for a
// single table. Everything specified for the WHERE clause is assumed to be an
// AND. This is mainly meant for basic INSERT, UPDATE, SELECT, DELETE
// statements on a table.
type TableSQL struct {
	Table              string   // the table from which to SELECT
	Columns            []string // the columns that will be SELECTed
	WhereColumns       []string // the where column names
	WhereComparisonOps []string // the comparison operator for the corresponding column index
	WhereConditions    []string // The conditional operator for Column pairs.
}
