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

// dbsql2go is a CLI tool to generate Go struct definitions from tables in a
// database.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mohae/dbsql2go"
	"github.com/mohae/dbsql2go/mysql"
)

const name = "dbsql2go"

var (
	dbType        string
	dbName        string
	pkgName       string
	server        string
	user          string
	password      string
	out           string
	filePerTable  bool
	includeDBName bool
)

func init() {
	flag.StringVar(&dbType, "rdbms", "", "the target RDBMS")
	flag.StringVar(&dbName, "db", "", "database name")
	flag.StringVar(&user, "user", "", "login user")
	flag.StringVar(&password, "password", "", "user's password")
	flag.StringVar(&server, "server", "", "server location")
	flag.StringVar(&pkgName, "package", "", "name of the package of which the generated code is a part; if empty,the database name will be used")
	flag.StringVar(&out, "out", "", "the output destination: if it doesn't end with a .go extension it will be assumed to be a path relative to the GOPATH/src dir. If empty, it will be the WD.")
	flag.BoolVar(&filePerTable, "separatefiles", false, "use a file per table; each file will use the table's name")
	flag.BoolVar(&includeDBName, "includedbname", false, "prefix each table file with the db name; only used with separate files")
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	flag.Parse()
	if dbType == "" {
		fmt.Fprintln(os.Stderr, "a rdbms must be specified")
		return 1
	}
	if dbName == "" {
		fmt.Fprintln(os.Stderr, "a db must be specified")
		return 1
	}
	if user == "" {
		fmt.Fprintln(os.Stderr, "a user must be specified")
		return 1
	}
	if password == "" {
		fmt.Fprintln(os.Stderr, "a password must be specified")
		return 1
	}
	// server is optional, depending on the db being used (for now). As such
	// it should be checked when applicable.

	// TODO: thinking about having it use the current dir as the package name
	// and adding a flag to use the dbname as the package name with that flag
	// taking precedence.
	if pkgName == "" {
		pkgName = dbName
	}
	// If the db request is not supported...
	typ, err := dbsql2go.ParseDBType(dbType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return 1
	}
	var DB dbsql2go.DBer
	var imp string // the db specific driver
	switch typ {
	case dbsql2go.MySQL:
		DB, err = mysql.New(server, user, password, dbName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s connect: %s\n", typ, err)
			return 1
		}
		imp = mysql.Import()
	}

	// Figure out the output destination.
	if out == "" {
		// get the wd
		out, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: get WD: %s\n", err)
			return 1
		}
		out = filepath.Join(out, dbName+".go")
	} else {
		// if out doesn't end with .go assume it is a dir within $GOPATH/src/
		if filepath.Ext(out) != ".go" {
			// make the fullpath first
			out = filepath.Join(os.Getenv("GOPATH"), out)
			// make sure the dir exists
			err = os.MkdirAll(out, 0766)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: directory: %s\n", err)
				return 1
			}
			out = filepath.Join(out, dbName+".go")
		}
	}

	// open the file
	f, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR|os.O_APPEND|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: open file: %s\n", err)
		return 1
	}
	defer f.Close()

	// TODO: add support for optional supported methods on generated structs???
	// Write the package name and imports to the file
	_, err = f.WriteString(fmt.Sprintf("//%s: %s struct definitions for %s\n", pkgName, typ, dbName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: package comment: %s\n", err)
		return 1
	}
	_, err = f.WriteString(fmt.Sprintf("package %s\n\nimport (\n\t\"database/sql\"\n\n\t%s\n)\n", pkgName, imp))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: package statements: %s\n", err)
		return 1
	}

	// dump all the Go table definitions to file
	// TODO: add support for file per table
	tables, err := DB.Tables()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: getting table info from %s: %s\n", dbName, err)
		return 1
	}

	for _, tbl := range tables {
		fmt.Println(tbl.Name())
		_, err = f.Write([]byte("\n\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: writing table separator lines: %s\n", err)
			return 1
		}

		b, err := tbl.GoFmt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: generating Go struct definition for %s.%s: %s\n", dbName, tbl.Name(), err)
			return 1
		}

		_, err = f.Write(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: write %s.%s definition to file: %s\n", dbName, tbl.Name(), err)
			return 1
		}
	}
	fmt.Printf("Go structs were generated from %s and written to %q\n", dbName, out)

	return 0
}