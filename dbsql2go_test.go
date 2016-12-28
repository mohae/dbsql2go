package dbsql2go

import "testing"

func TestParseDBType(t *testing.T) {
	tests := []struct {
		value string
		typ   DBType
		err   error
	}{
		{"", Unsupported, UnsupportedDBErr{Value: ""}},
		{"a", Unsupported, UnsupportedDBErr{Value: "a"}},
		{"MySQL", MySQL, nil},
		{"MYSQL", MySQL, nil},
		{"mysql", MySQL, nil},
		{"mYsQl", MySQL, nil},
		{"Postgres", Unsupported, UnsupportedDBErr{Value: "Postgres"}},
		{"postgres", Unsupported, UnsupportedDBErr{Value: "postgres"}},
		{"POSTGRES", Unsupported, UnsupportedDBErr{Value: "POSTGRES"}},
		{"pOsTgReS", Unsupported, UnsupportedDBErr{Value: "pOsTgReS"}},
		{"SQL Server", Unsupported, UnsupportedDBErr{Value: "SQL Server"}},
		{"sql server", Unsupported, UnsupportedDBErr{Value: "sql server"}},
		{"SQL SERVER", Unsupported, UnsupportedDBErr{Value: "SQL SERVER"}},
		{"sQl sERver", Unsupported, UnsupportedDBErr{Value: "sQl sERver"}},
		{"SQLServer", Unsupported, UnsupportedDBErr{Value: "SQLServer"}},
		{"sqlserver", Unsupported, UnsupportedDBErr{Value: "sqlserver"}},
		{"SQLSERVER", Unsupported, UnsupportedDBErr{Value: "SQLSERVER"}},
		{"sQlsERver", Unsupported, UnsupportedDBErr{Value: "sQlsERver"}},
		{"MSSQL Server", Unsupported, UnsupportedDBErr{Value: "MSSQL Server"}},
		{"mssql server", Unsupported, UnsupportedDBErr{Value: "mssql server"}},
		{"MSSQL SERVER", Unsupported, UnsupportedDBErr{Value: "MSSQL SERVER"}},
		{"msSql sERver", Unsupported, UnsupportedDBErr{Value: "msSql sERver"}},
		{"MSSQLServer", Unsupported, UnsupportedDBErr{Value: "MSSQLServer"}},
		{"mssqlserver", Unsupported, UnsupportedDBErr{Value: "mssqlserver"}},
		{"MSSQLSERVER", Unsupported, UnsupportedDBErr{Value: "MSSQLSERVER"}},
		{"msSqlsERver", Unsupported, UnsupportedDBErr{Value: "msSqlsERver"}},
		{"Oracle", Unsupported, UnsupportedDBErr{Value: "Oracle"}},
		{"oracle", Unsupported, UnsupportedDBErr{Value: "oracle"}},
		{"ORACLE", Unsupported, UnsupportedDBErr{Value: "ORACLE"}},
		{"OrAcLe", Unsupported, UnsupportedDBErr{Value: "OrAcLe"}},
		{"SQLite", Unsupported, UnsupportedDBErr{Value: "SQLite"}},
		{"sqlite", Unsupported, UnsupportedDBErr{Value: "sqlite"}},
		{"SQLITE", Unsupported, UnsupportedDBErr{Value: "SQLITE"}},
		{"SqLiTe", Unsupported, UnsupportedDBErr{Value: "SqLiTe"}},
	}
	for _, test := range tests {
		typ, err := ParseDBType(test.value)
		if err != test.err {
			t.Errorf("%s: got %v want %v", err, test.err)
			continue
		}
		if typ != test.typ {
			t.Errorf("%s: got %v want %v", typ, test.typ)
		}
	}
}
