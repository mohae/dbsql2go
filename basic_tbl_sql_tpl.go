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

import "html/template"

// the basic sql stuff for a single table go in this file.

var (
	SelectSQL *template.Template // Template to SELECT from a single table with only ANDs
	DeleteSQL *template.Template // Template to DELETE from a single table with only ANDs
	InsertSQL *template.Template // Template to INSERT from a single table with only ANDs
)

func init() {
	SelectSQL = template.Must(template.New("select").Parse(selectSQL))
	DeleteSQL = template.Must(template.New("delete").Parse(deleteSQL))
	InsertSQL = template.Must(template.New("insert").Parse(insertSQL))
}

// selectSQL is the template for selecting data from a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic SELECTs from a table.
var selectSQL = `SELECT
{{- range $i, $col := .Columns -}}
	{{- if eq $i 0 }} {{ $col -}}
	{{- else -}}
		, {{$col}}
	{{- end -}}
{{- end }}
FROM {{.Table}}
WHERE {{- range $i, $col := .Where -}}
	{{- if eq $i 0 }} {{ $col }} = ?
	{{- else }}
    AND {{ $col }} = ?
	{{- end -}}
{{- end -}}
`

// deleteSQL is the template for delecting data from a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic DELETEs from a table.
var deleteSQL = `DELETE FROM {{.Table}}
WHERE {{- range $i, $col := .Where -}}
	{{- if eq $i 0 }} {{ $col }} = ?
	{{- else }}
    AND {{ $col }} = ?
	{{- end -}}
{{- end -}}
`

// insertSQL is the template for inserting data from a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic INSERTs from a table.
var insertSQL = `INSERT INTO {{.Table}} ({{- range $i, $col := .Columns -}}
	{{- if eq $i 0 }}{{ $col -}}
	{{- else -}}
		, {{$col}}
	{{- end -}}
{{- end -}}
)
VALUES ({{- range $i, $col := .Columns -}}
{{- if eq $i 0 -}} ?
{{- else -}}
	, ?
{{- end -}}
{{- end -}}
)
WHERE {{- range $i, $col := .Where -}}
	{{- if eq $i 0 }} {{ $col }} = ?
	{{- else }}
    AND {{ $col }} = ?
	{{- end -}}
{{- end -}}
`
