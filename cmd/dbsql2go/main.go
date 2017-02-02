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
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohae/dbsql2go"
	"github.com/mohae/dbsql2go/mysql"
)

var exe = filepath.Base(os.Args[0]) // name of executable

var (
	dbType       string
	dbName       string
	pkgName      string
	server       string
	user         string
	password     string
	out          string
	filePerTable bool
)

func init() {
	flag.StringVar(&dbType, "rdbms", "", "the target RDBMS")
	flag.StringVar(&dbName, "db", "", "database name")
	flag.StringVar(&user, "user", "", "login user")
	flag.StringVar(&user, "u", "", "login 'user'")
	flag.StringVar(&password, "password", "", "user's password")
	flag.StringVar(&password, "p", "", "user's 'password'")
	flag.StringVar(&server, "server", "", "server location")
	flag.StringVar(&pkgName, "package", "", "name of the package of which the generated code is a part; if empty,the database name will be used")
	flag.StringVar(&out, "out", "", "the output destination: if it doesn't end with a .go extension it will be assumed to be a path relative to the GOPATH/src dir. If empty, it will be the WD.")
	flag.BoolVar(&filePerTable, "separatefiles", false, "use a file per table; each file will use the table's name")

	log.SetFlags(0)
	log.SetPrefix(exe + ": ")
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s Usage:\n", exe)
	fmt.Fprintf(os.Stderr, "  %s [FLAGS] -rdbms mysql -db dbname -user username -password password \n", exe)
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprint(os.Stderr, "Creates Go structs from a database.\n")
	fmt.Fprint(os.Stderr, "\n")
	fmt.Fprint(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	// take care of flag stuff first
	flag.Usage = usage
	flag.Parse()
	if flag.NFlag() < 4 {
		fmt.Fprintln(os.Stderr, "At least -rdbms, -db, -user (-u), and -password (-p) flags must be passed.\nAdditional flags may be required, depending on the target RDBMS.\n")
		flag.Usage()
		os.Exit(2)
	}
	if dbType == "" {
		log.Fatal("a rdbms must be specified")
	}
	if dbName == "" {
		log.Fatal("a db must be specified")
	}
	if user == "" {
		log.Fatal("a user must be specified")
	}
	if password == "" {
		log.Fatal("a password must be specified")
	}

	// TODO: thinking about having it use the current dir as the package name
	// and adding a flag to use the dbname as the package name with that flag
	// taking precedence.
	if pkgName == "" {
		pkgName = dbName
	}

	// If the db request is not supported...
	typ, err := dbsql2go.ParseDBType(dbType)
	if err != nil {
		log.Fatal("error: %s\n", err)
	}

	var DB dbsql2go.DBer
	var imp string // the db specific driver

	// Connect to the DB
	switch typ {
	case dbsql2go.MySQL:
		DB, err = mysql.New(server, user, password, dbName)
		if err != nil {
			log.Fatal("error: %s connect: %s\n", typ, err)
		}
		imp = mysql.Import()
	}

	// Get gets all of the information for the specified db.
	err = DB.Get()
	if err != nil {
		log.Fatal("%s: error: gathering of ", dbName)
	}

	// we don't defer close
	w, filename, err := setOutput()
	if err != nil {
		log.Fatal("error: %s", err)
	}

	// Now that the DB information has been gathered and the output set; process
	// the database info and write to the file.

	// TODO: add support for optional supported methods on generated structs???
	// Write the package name and imports to the file
	_, err = w.Write([]byte(fmt.Sprintf("//%s: %s struct definitions for database tables and views.\n// auto-generated by github.com/mohae/dbsql2go/cmd/dbsql2go\n", pkgName, typ)))
	if err != nil {
		w.(*os.File).Close()
		log.Fatal("error: package comment: %s\n", err)
	}

	if filePerTable {
		w.(*os.File).Close() // close the db file; the table specific ones will be written to their own.
	} else {
		_, err = writeTableFileComments(w, imp)
		if err != nil {
			w.(*os.File).Close()
			log.Fatal("error: package statements: %s\n", err)
		}
	}

	// Get all the tables; this also gathers all relevant db info.
	tables := DB.Tables()
	// dump all the Go table definitions to file
	// TODO: add support for file per table
	for _, tbl := range tables {
		if filePerTable {
			// open the file for this table
			w, err = os.OpenFile(filepath.Join(out, tbl.Name()+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766)
			if err != nil {
				log.Fatal("error: open file: %s\n", err)
			}

		}
		_, err = w.Write([]byte("\n\n"))
		if err != nil {
			log.Fatal("error: writing table separator lines: %s\n", err)
			w.(*os.File).Close()
		}

		err := tbl.GoFmt(w)
		if err != nil {
			w.(*os.File).Close()
			log.Fatal("error: generating Go struct definition for %s.%s: %s\n", dbName, tbl.Name(), err)
		}
		if filePerTable { // done writing to the file so close it.
			w.(*os.File).Close()
		}
	}

	if !filePerTable {
		out = filepath.Join(out, filename)
	}

	fmt.Printf("Go structs were generated from %s and written to %q\n", dbName, out)
	return
}

// setOutput sets the initial output writer, ensures the filePerTable flag is
// unset, if applicable, and sets the out destination properly.
//
// If out is empty, the default GOPATH (as of 1.8) will be used. If out
func setOutput() (w io.Writer, filename string, err error) {
	// Figure out the output destination. I feel bad about this hairball.
	if out == "stdout" {
		w = os.Stdin
		// for stdout filePerTable doesn't make sense; make sure its false
		filePerTable = false
		return w, "", nil
	}

	// declarations because of goto
	var (
		i   int
		gop string
	)

	if out == "" {
		// get the wd
		out, err = os.Getwd()
		if err != nil {
			return nil, "", err
		}
		goto output
	}

	// check out for $GOPATH and replace if it exists
	// TODO: support windows
	i = strings.Index(out, "$GOPATH")
	if i >= 0 {
		gop = os.Getenv("GOPATH")
		if gop == "" { // if it wasn't set, use Go's default path (1.8) + src
			gop = "$HOME/go/src"
		}

	}
	out = os.ExpandEnv(out) // expand any other env vars that may be in the path.

	// the out includes a filename, separate that out.
	if filepath.Ext(out) == ".go" {
		filename = filepath.Base(out) // set the filename in this instance
		out = filepath.Dir(out)
	}

output:
	if filePerTable && filename != "" {
		fmt.Printf("%s: the separatefile flag was set to true and the output destination include a .go extension, the output directory will be %q and each table will be written to its own file", exe, out)
	}

	// make sure the dir exists
	if out != "" {
		err = os.MkdirAll(out, 0766)
		if err != nil {
			return nil, "", fmt.Errorf("directory: %s", err)
		}
	}
	fmt.Printf("%s: writing generated .go files to the %q directory\n", exe, out)

	// open the file
	if filename == "" {
		filename = dbName + ".go"
	}

	w, err = os.OpenFile(filepath.Join(out, filename), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766)
	if err != nil {
		return nil, "", fmt.Errorf("open file: %s", err)
	}

	return w, filename, nil
}

func writeTableFileComments(w io.Writer, imp string) (n int, err error) {
	return w.Write([]byte(fmt.Sprintf("package %s\n\nimport (\n\t\"database/sql\"\n\n\t%s\n)\n", pkgName, imp)))
}
