package dbsql2go

import (
	"bytes"
	"testing"
)

var tables = []TableSQL{
	TableSQL{
		Columns:      []string{},
		Table:        "",
		WhereColumns: []string{},
	},
	TableSQL{
		Columns:      []string{"bar"},
		Table:        "",
		WhereColumns: []string{},
	},
	TableSQL{
		Columns:      []string{},
		Table:        "foo",
		WhereColumns: []string{},
	},
	TableSQL{
		Columns:      []string{},
		Table:        "",
		WhereColumns: []string{"id"},
	},
	TableSQL{
		Columns:      []string{"bar"},
		Table:        "",
		WhereColumns: []string{"id"},
	},
	TableSQL{
		Columns:      []string{},
		Table:        "foo",
		WhereColumns: []string{"id"},
	},
	TableSQL{
		Columns:      []string{"bar"},
		Table:        "foo",
		WhereColumns: []string{"id"},
	},
	TableSQL{
		Columns:      []string{"bar", "biz", "baz"},
		Table:        "foo",
		WhereColumns: []string{"id"},
	},
	TableSQL{
		Columns:      []string{"bar", "biz", "baz"},
		Table:        "foo",
		WhereColumns: []string{"id", "sid"},
	},
}

func TestTableSelectTemplate(t *testing.T) {
	expected := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"SELECT bar FROM foo WHERE id = ?",
		"SELECT bar, biz, baz FROM foo WHERE id = ?",
		"SELECT bar, biz, baz FROM foo WHERE id = ? AND sid = ?",
	}
	var buff bytes.Buffer
	for i, tbl := range tables {
		// it's just simpler to reset it here
		buff.Reset()
		err := SelectSQL.Execute(&buff, tbl)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != expected[i] {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want %x", i, buff.String(), expected[i])
			t.Errorf("%d: got %q want %q", i, buff.String(), expected[i])
		}
	}
}

func TestTableDELETETemplate(t *testing.T) {
	expected := []string{
		"",
		"",
		"DELETE FROM foo",
		"",
		"",
		"DELETE FROM foo WHERE id = ?",
		"DELETE FROM foo WHERE id = ?",
		"DELETE FROM foo WHERE id = ?",
		"DELETE FROM foo WHERE id = ? AND sid = ?",
	}
	var buff bytes.Buffer
	for i, tbl := range tables {
		// it's just simpler to reset it here
		buff.Reset()
		err := DeleteSQL.Execute(&buff, tbl)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != expected[i] {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want %x", i, buff.String(), expected[i])
			t.Errorf("%d: got %q want %q", i, buff.String(), expected[i])
		}
	}
}

func TestTableINSERTTemplate(t *testing.T) {
	expected := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"INSERT INTO foo (bar) VALUES (?)",
		"INSERT INTO foo (bar, biz, baz) VALUES (?, ?, ?)",
		"INSERT INTO foo (bar, biz, baz) VALUES (?, ?, ?)",
	}
	var buff bytes.Buffer
	for i, tbl := range tables {
		// it's just simpler to reset it here
		buff.Reset()
		err := InsertSQL.Execute(&buff, tbl)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != expected[i] {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want %x", i, buff.String(), expected[i])
			t.Errorf("%d: got %q want %q", i, buff.String(), expected[i])
		}
	}
}

func TestTableUpdateTemplate(t *testing.T) {
	expected := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"UPDATE foo SET bar = ? WHERE id = ?",
		"UPDATE foo SET bar = ?, biz = ?, baz = ? WHERE id = ?",
		"UPDATE foo SET bar = ?, biz = ?, baz = ? WHERE id = ? AND sid = ?",
	}
	var buff bytes.Buffer
	for i, tbl := range tables {
		// it's just simpler to reset it here
		buff.Reset()
		err := UpdateSQL.Execute(&buff, tbl)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != expected[i] {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want %x", i, buff.String(), expected[i])
			t.Errorf("%d: got %q want %q", i, buff.String(), expected[i])
		}
	}
}

var andOrtests = []struct {
	TableSQL
	expectedSQL   string
	expectedWhere string
}{
	{ //not enough comparison ops for the where columns == empty string
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id"},
			WhereComparisonOps: []string{">"},
			WhereConditions:    []string{},
		},
		"",
		"",
	},
	{ // no where condition == empty string
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id"},
			WhereComparisonOps: []string{">", "<"},
			WhereConditions:    []string{},
		},
		"",
		"",
	},
	{
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id"},
			WhereComparisonOps: []string{">", "<"},
			WhereConditions:    []string{"AND"},
		},
		`SELECT id, val FROM abc WHERE id > ? AND id < ?`,
		`WHERE id > arg[0] AND id < arg[1]`,
	},
	{
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id"},
			WhereComparisonOps: []string{"<", ">"},
			WhereConditions:    []string{"OR"},
		},
		`SELECT id, val FROM abc WHERE id < ? OR id > ?`,
		`WHERE id < arg[0] OR id > arg[1]`,
	},
	{
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id", "name"},
			WhereComparisonOps: []string{">", "<", "="},
			WhereConditions:    []string{"AND", "OR"},
		},
		`SELECT id, val FROM abc WHERE id > ? AND id < ? OR name = ?`,
		`WHERE id > arg[0] AND id < arg[1] OR name = arg[2]`,
	},
	{
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id", "val", "name"},
			WhereComparisonOps: []string{">", "<", "!=", "="},
			WhereConditions:    []string{"AND", "AND", "OR"},
		},
		`SELECT id, val FROM abc WHERE id > ? AND id < ? AND val != ? OR name = ?`,
		`WHERE id > arg[0] AND id < arg[1] AND val != arg[2] OR name = arg[3]`,
	},
	{
		TableSQL{
			Columns:            []string{"id", "val"},
			Table:              "abc",
			WhereColumns:       []string{"id", "id", "val", "val"},
			WhereComparisonOps: []string{">", "<", ">", "<"},
			WhereConditions:    []string{"AND", "AND", "AND"},
		},
		`SELECT id, val FROM abc WHERE id > ? AND id < ? AND val > ? AND val < ?`,
		`WHERE id > arg[0] AND id < arg[1] AND val > arg[2] AND val < arg[3]`,
	},
}

func TestTableSelectANDORSQLTemplate(t *testing.T) {

	var buff bytes.Buffer
	// first do the basic tests to see if the results are as expected; all
	// results should be empty string as the tables tests don't have the
	// necessary WhereComparisonOps and WhereConditions.
	for i, tbl := range tables {
		// it's just simpler to reset it here
		buff.Reset()
		err := SelectAndOrSQL.Execute(&buff, tbl)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != "" {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want \"\"", i, buff.String())
			t.Errorf("%d: got %q want \"\"", i, buff.String())
		}
	}

	// now do the actual tests
	for i, test := range andOrtests {
		// it's just simpler to reset it here
		buff.Reset()
		err := SelectAndOrSQL.Execute(&buff, test.TableSQL)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != test.expectedSQL {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want %x", i, buff.String(), test.expectedSQL)
			t.Errorf("%d: got %q want %q", i, buff.String(), test.expectedSQL)
		}
	}
}

func TestTableSelectANDORWhereCommentTemplate(t *testing.T) {

	var buff bytes.Buffer
	for i, test := range andOrtests {
		// it's just simpler to reset it here
		buff.Reset()
		err := SelectAndOrWhereComment.Execute(&buff, test.TableSQL)
		if err != nil {
			t.Errorf("%d: %s", i, err)
			continue
		}
		if buff.String() != test.expectedWhere {
			// use the hex values because it makes it easier to spot difference that
			// aren't always obvious visually, e.g. trailing blanks
			t.Errorf("%d: got %x want %x", i, buff.String(), test.expectedWhere)
			t.Errorf("%d: got %q want %q", i, buff.String(), test.expectedWhere)
		}
	}
}
