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

import "strings"

const (
	Unsupported DBType = iota
	MySQL
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
	PK       IndexType = iota + 1 // Primary Key
	FK                            // Foreign Key
	Unique                        // Unique
	Standard                      // An index that isn't one of the above.
)

//go:generate stringer -type=IndexType
type IndexType int

type UnsupportedDBErr struct {
	Value string
}

func (u UnsupportedDBErr) Error() string {
	return u.Value + " is not a supported database system"
}

type DBer interface {
	GetTables() error
	Tables() []Tabler
	GetIndexes() error
	Indexes() []Indexer
	GetViews() error
	Views() []Viewer
}

// Tabler
type Tabler interface {
	//Columns() []Column
	Name() string
	Schema() string
	Collation() string
	Go() ([]byte, error)
	GoFmt() ([]byte, error)
	//	SelectSQL() string
	//	InsertSQL() string
	//	DeleteSQL() string
}

// Indexer
type Indexer interface {
	Name() string // Just so that there's semething to fulfill until this gets fleshed out further.
}

// Index holds information about a given index.
type Index struct {
	Type     IndexType // type of Index
	Name     string    // Name of Index
	Table    string    // Index's table
	Cols     []string  // Index Columns, in order.
	RefTable string    // Referred to table for Foreign Keys
	RefCols  []string  // Referred to columns, in order, for Foreign Keys
}

// Viewer
type Viewer interface {
	Name() string // Just so that there's semething to fulfill until this gets fleshed out further.
}
