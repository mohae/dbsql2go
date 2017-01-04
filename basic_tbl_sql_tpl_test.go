package dbsql2go

import (
	"bytes"
	"testing"
)

var tables = []TableSQL{
	TableSQL{
		Columns: []string{},
		Table:   "",
		Where:   []string{},
	},
	TableSQL{
		Columns: []string{"bar"},
		Table:   "foo",
		Where:   []string{"id"},
	},
	TableSQL{
		Columns: []string{"bar", "biz", "baz"},
		Table:   "foo",
		Where:   []string{"id"},
	},
	TableSQL{
		Columns: []string{"bar", "biz", "baz"},
		Table:   "foo",
		Where:   []string{"id", "sid"},
	},
}

// while the behavior seems wrong; saves in atom results in in the elision of
// the blank, 0x20, from the first expected result which causes a test fail.
// This does not appear to be a gofmt issue as running gofmt does not result
// in the blank being elided. TODO: figure out what is happening
func TestTableSelectTemplate(t *testing.T) {
	expected := []string{
		`SELECT
FROM 
WHERE`,
		`SELECT bar
FROM foo
WHERE id == ?`,
		`SELECT bar, biz, baz
FROM foo
WHERE id == ?`,
		`SELECT bar, biz, baz
FROM foo
WHERE id == ?
    AND sid == ?`,
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
		}
	}
}

// while the behavior seems wrong; saves in atom results in in the elision of
// the blank, 0x20, from the first expected result which causes a test fail.
// This does not appear to be a gofmt issue as running gofmt does not result
// in the blank being elided. TODO: figure out what is happening
func TestTableDELETETemplate(t *testing.T) {
	expected := []string{
		`DELETE FROM 
WHERE`,
		`DELETE FROM foo
WHERE id == ?`,
		`DELETE FROM foo
WHERE id == ?`,
		`DELETE FROM foo
WHERE id == ?
    AND sid == ?`,
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
		}
	}
}
