// Copyright 2016 Joel Scoble.
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

// Table is a generalized table struct. Each supported database will embed
// this struct
type Tabler interface {
	Columns() []Column
	Name() string
	Tablespace() string
	Collation() string
}

// Column represents a column within a table.
type Column struct {
	Name     string      // Name of the column
	Datatype string      // Datatype of the column
	NotNull  bool        // If the column is NOT NULL
	Default  interface{} // Default value; if applicable
}
