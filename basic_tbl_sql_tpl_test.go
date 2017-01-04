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
		err := TableSelect.Execute(&buff, tbl)
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
