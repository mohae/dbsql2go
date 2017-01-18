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
			t.Errorf("%s: got %v want %v", test.value, err, test.err)
			continue
		}
		if typ != test.typ {
			t.Errorf("%s: got %v want %v", test.value, typ, test.typ)
		}
	}
}

func TestParseConstraintType(t *testing.T) {
	tests := []struct {
		value    string
		expected ConstraintType
		err      error
	}{
		{"PRIMARY KEY", PK, nil},
		{"PRIMARY", PK, nil},
		{"FOREIGN KEY", FK, nil},
		{"UNIQUE", Unique, nil},
		{"", UnknownConstraint, UnknownConstraintErr{""}},
		{"u", UnknownConstraint, UnknownConstraintErr{"u"}},
		{"alt", UnknownConstraint, UnknownConstraintErr{"alt"}},
		{"U", UnknownConstraint, UnknownConstraintErr{"U"}},
	}

	for _, test := range tests {
		typ, err := ParseConstraintType(test.value)
		if err != test.err {
			t.Errorf("%s: got %v want %v", test.value, err, test.err)
			continue
		}
		if typ != test.expected {
			t.Errorf("%s: got %v want %v", test.value, typ, test.expected)
		}
	}
}

func TestStringInComments(t *testing.T) {
	tests := []struct {
		line  string
		l     int
		lines []string
	}{
		{"", 10, nil},
		{"Hello", 10, []string{"// Hello"}},
		{"Hello World", 10, []string{"// Hello", "// World"}},
		{"This sentence is a meaningless one", 0, []string{"// This sentence is a meaningless one"}},
		{"This sentence is a meaningless one", 20, []string{"// This sentence is", "// a meaningless one"}},
		{"못 알아 듣겠어요 전혀 모르겠어요", 10, []string{"// 못 알아", "// 듣겠어요 전혀", "// 모르겠어요"}},
		// outlier, but if a word > l then use the whole word anyways
		{"hello Χαίρετε Здравствуйте", 10, []string{"// hello", "// Χαίρετε", "// Здравствуйте"}},
	}
	for i, test := range tests {
		lines := StringToComments(test.line, test.l)
		if len(lines) != len(test.lines) {
			t.Errorf("%d: got %d lines; want %d", i, len(lines), len(test.lines))
			t.Errorf("%s", lines)
			continue
		}
		for j, v := range lines {
			if v != test.lines[j] {
				t.Errorf("%d:%d: got %q; want %q", i, j, v, test.lines[j])
				continue
			}
		}
	}
}
