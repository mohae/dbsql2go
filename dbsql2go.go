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
package dbsql2go

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	Unsupported DBType = iota
	MySQL
)

const (
	commentStart = "// "
	LF           = []byte{'\n'}
)

//go:generate stringer -type=DBType
type DBType int

func ParseDBType(s string) (DBType, error) {
	v := strings.ToLower(s)
	switch v {
	case "mysql":
		return MySQL, nil
	default:
		return Unsupported, UnsupportedDBErr{Value: s}
	}
}

const (
	UnknownConstraint ConstraintType = iota
	PK
	FK
	Unique
)

//go:generate stringer -type=ConstraintType
// ConstarintType is the type of the table constraint.
type ConstraintType int

func ParseConstraintType(s string) (ConstraintType, error) {
	v := strings.ToLower(s)
	switch v {
	case "primary key", "primary":
		return PK, nil
	case "foreign key":
		return FK, nil
	case "unique":
		return Unique, nil
	default:
		return UnknownConstraint, UnknownConstraintErr{s}
	}
}

type UnknownConstraintErr struct {
	Value string
}

func (u UnknownConstraintErr) Error() string {
	return u.Value + " is not a known DB constraint type"
}

type UnsupportedDBErr struct {
	Value string
}

func (u UnsupportedDBErr) Error() string {
	return u.Value + " is not a supported database system"
}

type DBer interface {
	Get() error
	GetTables() error
	Tables() []Tabler
	GetIndexes() error
	GetConstraints() error
	GetViews() error
	Views() []Viewer
	UpdateTableConstraints() error
	UpdateTableIndexes()
}

// Tabler
type Tabler interface {
	//Columns() []Column
	Name() string
	Schema() string
	Collation() string
	Definition(io.Writer) error
	Go(io.Writer) error
	GoFmt(io.Writer) error
	ColumnNames() []string
	NonPKColumnNames() []string
	Indexes() []Index
	Constraints() []Constraint
	IsView() bool // If this is actually a view
	SelectSQLPK(io.Writer) error
	DeleteSQLPK(io.Writer) error
	InsertSQL(io.Writer) error
	GetPK() *Constraint
	StructName() string
	// SelectSQL
	//InsertSQL() string
	//	DeleteSQL() string
}

// Indexer
type Indexer interface {
	Name() string // Just so that there's semething to fulfill until this gets fleshed out further.
}

// Index holds information about a given index.
type Index struct {
	Type    string   // type of Index
	Primary bool     // if the Index is a primary key
	Name    string   // Name of Index
	Table   string   // Index's table
	Cols    []string // Index Columns, in order.
}

// Constraint holds information about a table's constraints, e.g. Primary Key.
type Constraint struct {
	Type      ConstraintType // The key or constraint type
	Name      string         // Name of key
	Table     string         // the table to which this key belongs.
	Cols      []string       // the columns that this key/constraint are on, in order.
	Fields    []string       // the Go struct field names corresponding to the table's column names.
	RefTable  string         // Referred to table for Foreign Keys
	RefCols   []string       // Referred to columns, in order, for Foreign Keys
	RefFields []string       // the Go struct field names corresponding to the table's column names.
}

// Viewer
type Viewer interface {
	Name() string // Just so that there's semething to fulfill until this gets fleshed out further.
}

// string to comments takes a string splits it up into comments l chars in
// length and returns a slice of comments. A comment starts with `// `. If
// l == 0; the comments will be 80 chars in length. The comment length refers
// the length in characters (runes) and not bytes.
func StringToComments(s string, l int) []string {
	if s == "" { // if the string is empty, no comment
		return nil
	}
	if l == 0 { // set the default if 0
		l = 80
	}
	// reduce the length by the length of commentBegin.
	l -= len(commentStart)
	var c []string
	if len(s) <= l {
		return append(c, fmt.Sprintf("%s%s", commentStart, s))
	}

	var (
		r []rune // line buffer
		b int    // current line length in characters
		k int    // k is the index of the last space in r
	)

	// separate out to words and comment lines
	for _, v := range s {
		if b > l {
			// only do if space was encountered
			if k != 0 {
				c = append(c, fmt.Sprintf("%s%s", commentStart, string(r[:k]))) // use everything up to the space
				if k+1 <= len(r) {
					r = r[k+1:] // if there were runes processed after the last seen space, keep them for next line
				} else {
					r = r[:0] // if the line happened on a space boundary, the next line starts out empty
				}
				b = len(r) // keep track of number of chars already in the next line
				k = 0      // reset space tracker
			}
		}
		if unicode.IsSpace(v) {
			k = b // set space index to current
		}
		r = append(r, v) // add the rune to the current line.
		b++              // increment the character count for the current line
	}
	c = append(c, fmt.Sprintf("%s%s", commentStart, string(r)))
	return c
}
