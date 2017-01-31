package dbsql2go

import "text/template"

// the basic sql stuff for a single table go in this file.

// Templates
var (
	SelectSQL *template.Template // Template to SELECT from a single table with only ANDs
	// SelectAndOrSQL allows for between and not between type statements on a single
	// table using the supplied conditional operators that correspond with the column
	// of the same index. BETWEEN is not used so that exclusive and inclusive can be
	// supported using: <, >, =, <=, >=. Since this template will generate SQL
	// that may return more than one row of data, these should be used in funcs that
	// return a slice of results; not as a method on a table.
	SelectAndOrSQL *template.Template
	DeleteSQL      *template.Template // Template to DELETE from a single table with only ANDs
	InsertSQL      *template.Template // Template to INSERT into a single table with only ANDs
	UpdateSQL      *template.Template // Template to UPDATE a row in a single table with only ANDs

	// Comment fragments
	SelectAndOrWhereComment *template.Template // The WHERE clause comment fragment for AndOR SQL funcs.
)

func init() {
	funcMap := template.FuncMap{
		"minusOne": minusOne,
	}
	SelectSQL = template.Must(template.New("select").Parse(selectSQL))
	SelectAndOrSQL = template.Must(template.New("selectandor").Funcs(funcMap).Parse(selectAndOrSQL))
	DeleteSQL = template.Must(template.New("delete").Parse(deleteSQL))
	InsertSQL = template.Must(template.New("insert").Parse(insertSQL))
	UpdateSQL = template.Must(template.New("update").Parse(updateSQL))

	SelectAndOrWhereComment = template.Must(template.New("selectandorcomment").Funcs(funcMap).Parse(selectAndOrWhereComment))
}

func minusOne(i int) int {
	return i - 1
}

// selectSQL is the template for selecting data from a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic SELECTs from a table.
var selectSQL = `{{ if and (ne .Table "") (gt (len .Columns) 0) -}}
SELECT
{{- range $i, $col := .Columns -}}
	{{- if eq $i 0 }} {{ $col -}}
	{{- else -}}
		, {{$col}}
	{{- end -}}
{{- end }} FROM {{.Table}}
{{- if gt (len .WhereColumns) 0 }} WHERE {{- range $i, $col := .WhereColumns -}}
	{{- if eq $i 0 }} {{ $col }} = ?
	{{- else }} AND {{ $col }} = ?
	{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

// selectAndOrSQL is the template for selecting data from a single table using
// multiple conditions and various conditional operators. This is primaraly
// meant for BETWEEN type operations using AND instead so that exclusive and
// inclusive evaluations can be done but this also can do more general
// comparisons limited by the limited logic of the WHERE clause generation
// in the template.
var selectAndOrSQL = `{{ $ComparisonMinus := minusOne (len .WhereComparisonOps) -}}
{{ if and (ne .Table "") (and (gt (len .Columns) 0) (and (gt (len .WhereColumns) 0) (and (eq (len .WhereColumns) (len .WhereComparisonOps)) (and (eq (len .WhereConditions) $ComparisonMinus))))) -}}
SELECT
{{- range $i, $col := .Columns -}}
	{{- if eq $i 0 }} {{ $col -}}
	{{- else -}}
		, {{$col}}
	{{- end -}}
{{- end }} FROM {{ .Table }} WHERE {{- range $i, $col := .WhereColumns -}}
	{{- if eq $i 0 }} {{ $col }} {{ index $.WhereComparisonOps $i }} ?
	{{- else }} {{index $.WhereConditions (minusOne ($i))}} {{ $col }} {{ index $.WhereComparisonOps $i }} ?
	{{- end -}}
{{- end -}}
{{- end -}}
`

// deleteSQL is the template for delecting data from a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic DELETEs from a table.
var deleteSQL = `{{ if ne .Table "" -}}
DELETE FROM {{.Table}}
{{- if gt (len .WhereColumns) 0 }} WHERE {{- range $i, $col := .WhereColumns -}}
	{{- if eq $i 0 }} {{ $col }} = ?
	{{- else }} AND {{ $col }} = ?
	{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

// insertSQL is the template for inserting data from a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic INSERTs into a table.
var insertSQL = `{{ if and (ne .Table "") (gt (len .Columns) 0) -}}
INSERT INTO {{.Table}} (
{{- range $i, $col := .Columns -}}
	{{- if eq $i 0 }}{{ $col -}}
	{{- else -}}, {{$col}}
	{{- end -}}
{{- end -}}
) VALUES ({{- range $i, $col := .Columns -}}
{{- if eq $i 0 -}} ?
{{- else }}, ?
{{- end -}}
{{- end -}}
)
{{- end -}}
`

// updateSQL is the template for updating a row in a single table. All
// columns in the WHERE field are assumed to use AND. Support for other
// conditions may be added in the future, but it complicates things, and,
// initially, this is meant to just create the basic UPDATES to a table row.
var updateSQL = `{{ if and (ne .Table "") (gt (len .Columns) 0) -}}
UPDATE {{.Table}} SET {{ range $i, $col := .Columns -}}
	{{- if eq $i 0 }}{{ $col }} = ?
	{{- else -}}, {{$col}} = ?
	{{- end -}}
{{- end -}}
{{ if gt (len .WhereColumns) 0 }} WHERE {{- range $i, $col := .WhereColumns -}}
	{{- if eq $i 0 }} {{ $col }} = ?
	{{- else }} AND {{ $col }} = ?
	{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

// selectAndOrWhereComment generates the example WHERE clause for the comments
// of a SELECT range func.
var selectAndOrWhereComment = `{{ $ComparisonMinus := minusOne (len .WhereComparisonOps) -}}
{{ if and (gt (len .WhereColumns) 0) (and (eq (len .WhereColumns) (len .WhereComparisonOps)) (and (eq (len .WhereConditions) $ComparisonMinus))) -}}
WHERE {{- range $i, $col := .WhereColumns -}}
	{{- if eq $i 0 }} {{ $col }} {{ index $.WhereComparisonOps $i }} arg[{{ $i }}]
	{{- else }} {{index $.WhereConditions (minusOne ($i))}} {{ $col }} {{ index $.WhereComparisonOps $i }} arg[{{ $i }}]
	{{- end -}}
{{- end -}}
{{- end -}}
`
