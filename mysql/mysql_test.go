package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/mohae/dbsql2go"
)

var (
	server         = "localhost"
	user           = "testuser"
	password       = "testuser"
	testDB         = "dbsql_test"
	testTablespace = "dbsql_tablespace"
)

var createTables = []string{
	`CREATE TABLE abc (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code CHAR(12) UNIQUE NOT NULL,
		description VARCHAR(20) NOT NULL,
		tiny TINYINT DEFAULT 3,
		small SMALLINT DEFAULT 11,
		medium MEDIUMINT DEFAULT 42,
		ger INTEGER,
		big BIGINT,
		cost DECIMAL,
		created TIMESTAMP
	)
	CHARACTER SET latin1 COLLATE latin1_swedish_ci`,
	`CREATE TABLE abc_nn (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code CHAR(12) UNIQUE NOT NULL,
		description VARCHAR(20) NOT NULL,
		tiny TINYINT NOT NULL,
		small SMALLINT NOT NULL,
		medium MEDIUMINT NOT NULL,
		ger INTEGER NOT NULL,
		big BIGINT NOT NULL,
		cost DECIMAL NOT NULL,
		created TIMESTAMP
	)
	CHARACTER SET latin1 COLLATE latin1_swedish_ci`,
	`CREATE TABLE def (
		id INT AUTO_INCREMENT PRIMARY KEY,
		d_date DATE,
		d_datetime DATETIME,
		d_time TIME,
		d_year YEAR,
		size ENUM('small', 'medium', 'large'),
		a_set SET('a', 'b', 'c'),
		INDEX (id, d_datetime)
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE def_nn (
		id INT AUTO_INCREMENT PRIMARY KEY,
		d_date DATE NOT NULL,
		d_datetime DATETIME NOT NULL,
		d_time TIME NOT NULL,
		d_year YEAR NOT NULL,
		size ENUM('small', 'med', 'large') NOT NULL,
		a_set SET('1', '2', '3') NOT NULL,
		INDEX (id, d_datetime)
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE ghi (
		id INT,
		val INT,
		def_id INT,
		def_datetime DATETIME,
		tiny_stuff TINYBLOB,
		stuff BLOB,
		med_stuff MEDIUMBLOB,
		long_stuff LONGBLOB,
		INDEX (val),
		FOREIGN KEY fk_def(def_id, def_datetime) REFERENCES def(id, d_datetime)
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE ghi_nn (
		id INT NOT NULL,
		val INT NOT NULL,
		def_id INT NOT NULL,
		def_datetime DATETIME NOT NULL,
		tiny_stuff TINYBLOB NOT NULL,
		stuff BLOB NOT NULL,
		med_stuff MEDIUMBLOB NOT NULL,
		long_stuff LONGBLOB NOT NULL,
		INDEX (val),
		FOREIGN KEY fk_def(def_id, def_datetime) REFERENCES def_nn(id, d_datetime)
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE jkl (
		id INT AUTO_INCREMENT,
		fid INT,
		tiny_txt TINYTEXT,
		txt TEXT,
		med_txt MEDIUMTEXT,
		long_txt LONGTEXT,
		bin BINARY(3),
		var_bin VARBINARY(12),
		PRIMARY KEY (id, fid),
		INDEX(fid),
		FOREIGN KEY(fid) REFERENCES def(id)
		ON UPDATE CASCADE
		ON DELETE RESTRICT
	)
	CHARACTER SET ascii COLLATE ascii_general_ci`,
	`CREATE TABLE jkl_nn (
		id INT AUTO_INCREMENT,
		fid INT,
		tiny_txt TINYTEXT NOT NULL,
		txt TEXT NOT NULL,
		med_txt MEDIUMTEXT NOT NULL,
		long_txt LONGTEXT NOT NULL,
		bin BINARY(3) NOT NULL,
		var_bin VARBINARY(12) NOT NULL,
		PRIMARY KEY (id, fid),
		INDEX(fid),
		FOREIGN KEY(fid) REFERENCES def(id)
		ON UPDATE CASCADE
		ON DELETE RESTRICT
	)
	CHARACTER SET ascii COLLATE ascii_general_ci`,
	`CREATE TABLE mno (
		id INT AUTO_INCREMENT PRIMARY KEY,
		geo GEOMETRY,
		pt POINT,
		lstring LINESTRING,
		poly POLYGON,
		multi_pt MULTIPOINT,
		multi_lstring MULTILINESTRING,
		multi_polygon MULTIPOLYGON,
		geo_collection GEOMETRYCOLLECTION
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE mno_nn (
		id INT AUTO_INCREMENT PRIMARY KEY,
		geo GEOMETRY NOT NULL,
		pt POINT NOT NULL,
		lstring LINESTRING NOT NULL,
		poly POLYGON NOT NULL,
		multi_pt MULTIPOINT NOT NULL,
		multi_lstring MULTILINESTRING NOT NULL,
		multi_polygon MULTIPOLYGON NOT NULL,
		geo_collection GEOMETRYCOLLECTION NOT NULL
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	//	`CREATE TABLE pqr (
	//		id INT AUTO_INCREMENT PRIMARY KEY,
	//		jsn JSON DEFAULT NULL
	//	)`,
}

var createViews = []string{
	`CREATE OR REPLACE VIEW abc_v
	AS SELECT id, code, description
	FROM abc
	ORDER by code`,
	`CREATE OR REPLACE VIEW defghi_v
	AS SELECT a.id AS aid, b.id as bid, a.d_datetime, a.size, b.stuff
	FROM def AS a, ghi AS b
	WHERE a.id = b.def_id
	ORDER by a.id, a.size, b.def_id`,
}

var tableDefs = []Table{
	Table{ // 0
		name: "abc", r: 'a', structName: "Abc", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "code", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "char", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "char(12)",
				Key: "UNI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Code",
			},
			Column{
				Name: "description", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varchar", CharMaxLen: sql.NullInt64{Int64: 20, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 20, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "varchar(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Description",
			},
			Column{
				Name: "tiny", OrdinalPosition: 4, Default: sql.NullString{String: "3", Valid: true},
				IsNullable: "YES", DataType: "tinyint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 3, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyint(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Tiny",
			},
			Column{
				Name: "small", OrdinalPosition: 5, Default: sql.NullString{String: "11", Valid: true},
				IsNullable: "YES", DataType: "smallint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 5, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "smallint(6)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Small",
			},
			Column{
				Name: "medium", OrdinalPosition: 6, Default: sql.NullString{String: "42", Valid: true},
				IsNullable: "YES", DataType: "mediumint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 7, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumint(9)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Medium",
			},
			Column{
				Name: "ger", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Ger",
			},
			Column{
				Name: "big", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "bigint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 19, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "bigint(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Big",
			},
			Column{
				Name: "cost", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "decimal", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "decimal(10,0)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Cost",
			},
			Column{
				Name: "created", OrdinalPosition: 10, Default: sql.NullString{String: "CURRENT_TIMESTAMP", Valid: true},
				IsNullable: "NO", DataType: "timestamp", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "timestamp",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Created",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "code", Table: "abc", Cols: []string{"code"}},
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "abc", Cols: []string{"id"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.Unique, Name: "code", Table: "abc",
				Cols: []string{"code"}, Fields: []string{"Code"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "abc",
				Cols: []string{"id"}, Fields: []string{"ID"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 1
		name: "abc_nn", r: 'a', structName: "AbcNn", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "code", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "char", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "char(12)",
				Key: "UNI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Code",
			},
			Column{
				Name: "description", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varchar", CharMaxLen: sql.NullInt64{Int64: 20, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 20, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "varchar(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Description",
			},
			Column{
				Name: "tiny", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinyint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 3, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyint(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Tiny",
			},
			Column{
				Name: "small", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "smallint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 5, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "smallint(6)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Small",
			},
			Column{
				Name: "medium", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 7, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumint(9)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Medium",
			},
			Column{
				Name: "ger", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Ger",
			},
			Column{
				Name: "big", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "bigint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 19, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "bigint(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Big",
			},
			Column{
				Name: "cost", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "decimal", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "decimal(10,0)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Cost",
			},
			Column{
				Name: "created", OrdinalPosition: 10, Default: sql.NullString{String: "CURRENT_TIMESTAMP", Valid: true},
				IsNullable: "NO", DataType: "timestamp", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "timestamp",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Created",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "code", Table: "abc_nn", Cols: []string{"code"}},
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "abc_nn", Cols: []string{"id"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.Unique, Name: "code", Table: "abc_nn",
				Cols: []string{"code"}, Fields: []string{"Code"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "abc_nn",
				Cols: []string{"id"}, Fields: []string{"ID"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 2
		name: "abc_v", r: 'a', structName: "AbcV", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "code", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "char", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "char(12)",
				Key: "UNI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Code",
			},
			Column{
				Name: "description", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varchar", CharMaxLen: sql.NullInt64{Int64: 20, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 20, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "varchar(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Description",
			},
		},
		Typ: "VIEW", Engine: sql.NullString{String: "", Valid: false},
		collation: sql.NullString{String: "", Valid: false}, Comment: "VIEW",
	},
	Table{ // 3
		name: "def", r: 'd', structName: "Def", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "d_date", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "date", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "date",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DDate",
			},
			Column{
				Name: "d_datetime", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DDatetime",
			},
			Column{
				Name: "d_time", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "time", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "time",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DTime",
			},
			Column{
				Name: "d_year", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "year", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "year(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DYear",
			},
			Column{
				Name: "size", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "enum", CharMaxLen: sql.NullInt64{Int64: 6, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 18, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "enum('small','medium','large')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Size",
			},
			Column{
				Name: "a_set", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "set", CharMaxLen: sql.NullInt64{Int64: 5, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 15, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "set('a','b','c')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ASet",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "id", Table: "def", Cols: []string{"id", "d_datetime"}},
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "def", Cols: []string{"id"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "def",
				Cols: []string{"id"}, Fields: []string{"ID"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 4
		name: "def_nn", r: 'd', structName: "DefNn", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "d_date", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "date", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "date",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DDate",
			},
			Column{
				Name: "d_datetime", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DDatetime",
			},
			Column{
				Name: "d_time", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "time", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "time",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DTime",
			},
			Column{
				Name: "d_year", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "year", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "year(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DYear",
			},
			Column{
				Name: "size", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "enum", CharMaxLen: sql.NullInt64{Int64: 5, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 15, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "enum('small','medium','large')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Size",
			},
			Column{
				Name: "a_set", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "set", CharMaxLen: sql.NullInt64{Int64: 5, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 15, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "set('a','b','c')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ASet",
			},
		}, // 5
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "id", Table: "def_nn", Cols: []string{"id", "d_datetime"}},
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "def_nn", Cols: []string{"id"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "def_nn",
				Cols: []string{"id"}, Fields: []string{"ID"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 6
		name: "defghi_v", r: 'd', structName: "DefghiV", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "aid", OrdinalPosition: 1, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Aid",
			},
			Column{
				Name: "bid", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Bid",
			},
			Column{
				Name: "d_datetime", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DDatetime",
			},
			Column{
				Name: "size", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "enum", CharMaxLen: sql.NullInt64{Int64: 6, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 18, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "enum('small','medium','large')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Size",
			},
			Column{
				Name: "stuff", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Stuff",
			},
		},
		Typ: "VIEW", Engine: sql.NullString{String: "", Valid: false},
		collation: sql.NullString{String: "", Valid: false}, Comment: "VIEW",
	},
	Table{ // 7
		name: "ghi", r: 'g', structName: "Ghi", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "val", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Val",
			},
			Column{
				Name: "def_id", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DefID",
			},
			Column{
				Name: "def_datetime", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DefDatetime",
			},
			Column{
				Name: "tiny_stuff", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "tinyblob", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "TinyStuff",
			},
			Column{
				Name: "stuff", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Stuff",
			},
			Column{
				Name: "med_stuff", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "mediumblob", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MedStuff",
			},
			Column{
				Name: "long_stuff", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "longblob", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "longblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "LongStuff",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "fk_def", Table: "ghi", Cols: []string{"def_id", "def_datetime"}},
			{Type: "BTREE", Primary: false, Name: "val", Table: "ghi", Cols: []string{"val"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.FK, Name: "ghi_ibfk_1", Table: "ghi",
				Cols: []string{"def_id", "def_datetime"}, Fields: []string{"DefID", "DefDatetime"}, RefTable: "def",
				RefCols: []string{"id", "d_datetime"}, RefFields: []string{"ID", "DDatetime"},
			},
		},
	},
	Table{ // 8
		name: "ghi_nn", r: 'g', structName: "GhiNn", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "val", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Val",
			},
			Column{
				Name: "def_id", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DefID",
			},
			Column{
				Name: "def_datetime", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "DefDatetime",
			},
			Column{
				Name: "tiny_stuff", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinyblob", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "TinyStuff",
			},
			Column{
				Name: "stuff", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Stuff",
			},
			Column{
				Name: "med_stuff", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumblob", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MedStuff",
			},
			Column{
				Name: "long_stuff", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "longblob", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "longblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "LongStuff",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "fk_def", Table: "ghi_nn", Cols: []string{"def_id", "def_datetime"}},
			{Type: "BTREE", Primary: false, Name: "val", Table: "ghi_nn", Cols: []string{"val"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.FK, Name: "ghi_nn_ibfk_1", Table: "ghi_nn",
				Cols: []string{"def_id", "def_datetime"}, Fields: []string{"DefID", "DefDatetime"}, RefTable: "def_nn",
				RefCols: []string{"id", "d_datetime"}, RefFields: []string{"ID", "DDatetime"},
			},
		},
	},
	Table{ // 9
		name: "jkl", r: 'j', structName: "Jkl", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "fid", OrdinalPosition: 2, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "FID",
			},
			Column{
				Name: "tiny_txt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "tinytext", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "tinytext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "TinyTxt",
			},
			Column{
				Name: "txt", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "text", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "text",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Txt",
			},
			Column{
				Name: "med_txt", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "mediumtext", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "mediumtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MedTxt",
			},
			Column{
				Name: "long_txt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "longtext", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "longtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "LongTxt",
			},
			Column{
				Name: "bin", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "binary", CharMaxLen: sql.NullInt64{Int64: 3, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 3, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "binary(3)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Bin",
			},
			Column{
				Name: "var_bin", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "varbinary", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "varbinary(12)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "VarBin",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "fid", Table: "jkl", Cols: []string{"fid"}},
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "jkl", Cols: []string{"id", "fid"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.FK, Name: "jkl_ibfk_1", Table: "jkl",
				Cols: []string{"fid"}, Fields: []string{"Fid"}, RefTable: "def",
				RefCols: []string{"id"}, RefFields: []string{"ID"},
			},
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "jkl",
				Cols: []string{"id", "fid"}, Fields: []string{"ID", "Fid"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 10
		name: "jkl_nn", r: 'j', structName: "JklNn", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "fid", OrdinalPosition: 2, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "FID",
			},
			Column{
				Name: "tiny_txt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinytext", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "tinytext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "TinyTxt",
			},
			Column{
				Name: "txt", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "text", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "text",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Txt",
			},
			Column{
				Name: "med_txt", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumtext", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "mediumtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MedTxt",
			},
			Column{
				Name: "long_txt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "longtext", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "longtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "LongTxt",
			},
			Column{
				Name: "bin", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "binary", CharMaxLen: sql.NullInt64{Int64: 3, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 3, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "binary(3)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Bin",
			},
			Column{
				Name: "var_bin", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varbinary", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "varbinary(12)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "VarBin",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: false, Name: "fid", Table: "jkl_nn", Cols: []string{"fid"}},
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "jkl_nn", Cols: []string{"id", "fid"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.FK, Name: "jkl_nn_ibfk_1", Table: "jkl_nn",
				Cols: []string{"fid"}, Fields: []string{"Fid"}, RefTable: "def",
				RefCols: []string{"id"}, RefFields: []string{"ID"},
			},
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "jkl_nn",
				Cols: []string{"id", "fid"}, Fields: []string{"ID", "Fid"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 11
		name: "mno", r: 'm', structName: "Mno", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "geo", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "geometry", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometry",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Geo",
			},
			Column{
				Name: "pt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "point", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "point",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Pt",
			},
			Column{
				Name: "lstring", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "linestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "linestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Lstring",
			},
			Column{
				Name: "poly", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "polygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "polygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Poly",
			},
			Column{
				Name: "multi_pt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "multipoint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipoint",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MultiPt",
			},
			Column{
				Name: "multi_lstring", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "multilinestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multilinestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MultiLstring",
			},
			Column{
				Name: "multi_polygon", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "multipolygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipolygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MultiPolygon",
			},
			Column{
				Name: "geo_collection", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "geometrycollection", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometrycollection",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "GeoCollection",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "mno", Cols: []string{"id"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "mno",
				Cols: []string{"id"}, Fields: []string{"ID"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
	Table{ // 12
		name: "mno_nn", r: 'm', structName: "MnoNn", schema: "dbsql_test",
		columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "ID",
			},
			Column{
				Name: "geo", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "geometry", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometry",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Geo",
			},
			Column{
				Name: "pt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "point", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "point",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Pt",
			},
			Column{
				Name: "lstring", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "linestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "linestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Lstring",
			},
			Column{
				Name: "poly", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "polygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "polygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "Poly",
			},
			Column{
				Name: "multi_pt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "multipoint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipoint",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MultiPt",
			},
			Column{
				Name: "multi_lstring", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "multilinestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multilinestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MultiLstring",
			},
			Column{
				Name: "multi_polygon", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "multipolygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipolygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "MultiPolygon",
			},
			Column{
				Name: "geo_collection", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "geometrycollection", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometrycollection",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "", fieldName: "GeoCollection",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
		indexes: []dbsql2go.Index{
			{Type: "BTREE", Primary: true, Name: "PRIMARY", Table: "mno_nn", Cols: []string{"id"}},
		},
		constraints: []dbsql2go.Constraint{
			{
				Type: dbsql2go.PK, Name: "PRIMARY", Table: "mno_nn",
				Cols: []string{"id"}, Fields: []string{"ID"}, RefTable: "",
				RefCols: nil, RefFields: nil,
			},
		},
	},
}

var tableDefsString = []string{
	`type Abc struct {
	ID int32
	Code string
	Description string
	Tiny sql.NullInt64
	Small sql.NullInt64
	Medium sql.NullInt64
	Ger sql.NullInt64
	Big sql.NullInt64
	Cost sql.NullFloat64
	Created mysql.NullTime
}
`,
	`type AbcNn struct {
	ID int32
	Code string
	Description string
	Tiny int8
	Small int16
	Medium int32
	Ger int32
	Big int64
	Cost float64
	Created mysql.NullTime
}
`,
	`type AbcV struct {
	ID int32
	Code string
	Description string
}
`,
	`type Def struct {
	ID int32
	DDate mysql.NullTime
	DDatetime mysql.NullTime
	DTime sql.NullString
	DYear sql.NullString
	Size sql.NullString
	ASet sql.NullString
}
`,
	`type DefNn struct {
	ID int32
	DDate mysql.NullTime
	DDatetime mysql.NullTime
	DTime string
	DYear string
	Size string
	ASet string
}
`,
	`type DefghiV struct {
	Aid int32
	Bid sql.NullInt64
	DDatetime mysql.NullTime
	Size sql.NullString
	Stuff []byte
}
`,
	`type Ghi struct {
	ID sql.NullInt64
	Val sql.NullInt64
	DefID sql.NullInt64
	DefDatetime mysql.NullTime
	TinyStuff []byte
	Stuff []byte
	MedStuff []byte
	LongStuff []byte
}
`,
	`type GhiNn struct {
	ID int32
	Val int32
	DefID int32
	DefDatetime mysql.NullTime
	TinyStuff []byte
	Stuff []byte
	MedStuff []byte
	LongStuff []byte
}
`,
	`type Jkl struct {
	ID int32
	Fid sql.NullInt64
	TinyTxt []byte
	Txt []byte
	MedTxt []byte
	LongTxt []byte
	Bin []byte
	VarBin []byte
}
`,
	`type JklNn struct {
	ID int32
	Fid int32
	TinyTxt []byte
	Txt []byte
	MedTxt []byte
	LongTxt []byte
	Bin []byte
	VarBin []byte
}
`,
}

var structDefs = []string{
	`type Abc struct {
	ID          int32
	Code        string
	Description string
	Tiny        sql.NullInt64
	Small       sql.NullInt64
	Medium      sql.NullInt64
	Ger         sql.NullInt64
	Big         sql.NullInt64
	Cost        sql.NullFloat64
	Created     mysql.NullTime
}

func (a *Abc) Select(db *sql.DB) error {
	err := db.QueryRow("SELECT id, code, description, tiny, small, medium, ger, big, cost, created FROM abc WHERE id = ?", a.ID).Scan(&a.ID, &a.Code, &a.Description, &a.Tiny, &a.Small, &a.Medium, &a.Ger, &a.Big, &a.Cost, &a.Created)
	if err != nil {
		return err
	}
	return nil
}

func (a *Abc) Delete(db *sql.DB) (n int, err error) {
	res, err := db.Exec("DELETE FROM abc WHERE id = ?", a.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
`,
	`type AbcNn struct {
	ID          int32
	Code        string
	Description string
	Tiny        int8
	Small       int16
	Medium      int32
	Ger         int32
	Big         int64
	Cost        float64
	Created     mysql.NullTime
}

func (a *AbcNn) Select(db *sql.DB) error {
	err := db.QueryRow("SELECT id, code, description, tiny, small, medium, ger, big, cost, created FROM abc_nn WHERE id = ?", a.ID).Scan(&a.ID, &a.Code, &a.Description, &a.Tiny, &a.Small, &a.Medium, &a.Ger, &a.Big, &a.Cost, &a.Created)
	if err != nil {
		return err
	}
	return nil
}

func (a *AbcNn) Delete(db *sql.DB) (n int, err error) {
	res, err := db.Exec("DELETE FROM abc_nn WHERE id = ?", a.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
`,
	`type AbcV struct {
	ID          int32
	Code        string
	Description string
}
`,
	`type Def struct {
	ID        int32
	DDate     mysql.NullTime
	DDatetime mysql.NullTime
	DTime     sql.NullString
	DYear     sql.NullString
	Size      sql.NullString
	ASet      sql.NullString
}

func (d *Def) Select(db *sql.DB) error {
	err := db.QueryRow("SELECT id, d_date, d_datetime, d_time, d_year, size, a_set FROM def WHERE id = ?", d.ID).Scan(&d.ID, &d.DDate, &d.DDatetime, &d.DTime, &d.DYear, &d.Size, &d.ASet)
	if err != nil {
		return err
	}
	return nil
}

func (d *Def) Delete(db *sql.DB) (n int, err error) {
	res, err := db.Exec("DELETE FROM def WHERE id = ?", d.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
`,
	`type DefNn struct {
	ID        int32
	DDate     mysql.NullTime
	DDatetime mysql.NullTime
	DTime     string
	DYear     string
	Size      string
	ASet      string
}

func (d *DefNn) Select(db *sql.DB) error {
	err := db.QueryRow("SELECT id, d_date, d_datetime, d_time, d_year, size, a_set FROM def_nn WHERE id = ?", d.ID).Scan(&d.ID, &d.DDate, &d.DDatetime, &d.DTime, &d.DYear, &d.Size, &d.ASet)
	if err != nil {
		return err
	}
	return nil
}

func (d *DefNn) Delete(db *sql.DB) (n int, err error) {
	res, err := db.Exec("DELETE FROM def_nn WHERE id = ?", d.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
`,
	`type DefghiV struct {
	Aid       int32
	Bid       sql.NullInt64
	DDatetime mysql.NullTime
	Size      sql.NullString
	Stuff     []byte
}
`,
	`type Ghi struct {
	ID          sql.NullInt64
	Val         sql.NullInt64
	DefID       sql.NullInt64
	DefDatetime mysql.NullTime
	TinyStuff   []byte
	Stuff       []byte
	MedStuff    []byte
	LongStuff   []byte
}
`,
	`type GhiNn struct {
	ID        int32
	Val       int32
	DefID     int32
	DDatetime mysql.NullTime
	TinyStuff []byte
	Stuff     []byte
	MedStuff  []byte
	LongStuff []byte
}
`,
	`type Jkl struct {
	ID      int32
	Fid     sql.NullInt64
	TinyTxt []byte
	Txt     []byte
	MedTxt  []byte
	LongTxt []byte
	Bin     []byte
	VarBin  []byte
}

func (j *Jkl) Select(db *sql.DB) error {
	err := db.QueryRow("SELECT id, fid, tiny_txt, txt, med_txt, long_txt, bin, var_bin FROM jkl WHERE id = ?", j.ID).Scan(&j.ID, &j.Fid, &j.TinyTxt, &j.Txt, &j,MedTxt, &j.LongTxt, &j.Bin, &j.VarBin)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jkl) Delete(db *sql.DB) (n int, err error) {
	res, err := db.Exec("DELETE FROM jkl WHERE id = ?", j.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
`,
	`type JklNn struct {
	ID      int32
	Fid     int32
	TinyTxt []byte
	Txt     []byte
	MedTxt  []byte
	LongTxt []byte
	Bin     []byte
	VarBin  []byte
}

func (j *JklNn) Select(db *sql.DB) error {
	err := db.QueryRow("SELECT id, fid, tiny_txt, txt, med_txt, long_txt, bin, var_bin FROM jkl_nn WHERE id = ?", j.ID).Scan(&j.ID, &j.Fid, &j.TinyTxt, &j.Txt, &j,MedTxt, &j.LongTxt, &j.Bin, &j.VarBin)
	if err != nil {
		return err
	}
	return nil
}

func (j *JklNn) Delete(db *sql.DB) (n int, err error) {
	res, err := db.Exec("DELETE FROM jkl_nn WHERE id = ?", j.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
`,
}

var indexes = []Index{
	{
		Table: "abc", NonUnique: 0, Schema: "dbsql_test", name: "code",
		SeqInIndex: 1, Column: "code", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "abc", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "abc_nn", NonUnique: 0, Schema: "dbsql_test", name: "code",
		SeqInIndex: 1, Column: "code", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "abc_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "def", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "def", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 2, Column: "d_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "def", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "def_nn", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "def_nn", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 2, Column: "d_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "def_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "ghi", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 1, Column: "def_id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "ghi", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 2, Column: "def_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "ghi", NonUnique: 1, Schema: "dbsql_test", name: "val",
		SeqInIndex: 1, Column: "val", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "ghi_nn", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 1, Column: "def_id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "ghi_nn", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 2, Column: "def_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "ghi_nn", NonUnique: 1, Schema: "dbsql_test", name: "val",
		SeqInIndex: 1, Column: "val", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "jkl", NonUnique: 1, Schema: "dbsql_test", name: "fid",
		SeqInIndex: 1, Column: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "jkl", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "jkl", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 2, Column: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "jkl_nn", NonUnique: 1, Schema: "dbsql_test", name: "fid",
		SeqInIndex: 1, Column: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "jkl_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "jkl_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 2, Column: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "mno", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		Table: "mno_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, Column: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
}

var constraints = []Constraint{
	{"code", "UNIQUE", "abc", "code", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "abc", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"code", "UNIQUE", "abc_nn", "code", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "abc_nn", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "def", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "def_nn", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"ghi_ibfk_1", "FOREIGN KEY", "ghi", "def_id", 1, sql.NullInt64{Int64: 1, Valid: true}, sql.NullString{String: "def", Valid: true}, sql.NullString{String: "id", Valid: true}},
	{"ghi_ibfk_1", "FOREIGN KEY", "ghi", "def_datetime", 2, sql.NullInt64{Int64: 2, Valid: true}, sql.NullString{String: "def", Valid: true}, sql.NullString{String: "d_datetime", Valid: true}},
	{"ghi_nn_ibfk_1", "FOREIGN KEY", "ghi_nn", "def_id", 1, sql.NullInt64{Int64: 1, Valid: true}, sql.NullString{String: "def_nn", Valid: true}, sql.NullString{String: "id", Valid: true}},
	{"ghi_nn_ibfk_1", "FOREIGN KEY", "ghi_nn", "def_datetime", 2, sql.NullInt64{Int64: 2, Valid: true}, sql.NullString{String: "def_nn", Valid: true}, sql.NullString{String: "d_datetime", Valid: true}},
	{"jkl_ibfk_1", "FOREIGN KEY", "jkl", "fid", 1, sql.NullInt64{Int64: 1, Valid: true}, sql.NullString{String: "def", Valid: true}, sql.NullString{String: "id", Valid: true}},
	{"PRIMARY", "PRIMARY KEY", "jkl", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "jkl", "fid", 2, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"jkl_nn_ibfk_1", "FOREIGN KEY", "jkl_nn", "fid", 1, sql.NullInt64{Int64: 1, Valid: true}, sql.NullString{String: "def", Valid: true}, sql.NullString{String: "id", Valid: true}},
	{"PRIMARY", "PRIMARY KEY", "jkl_nn", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "jkl_nn", "fid", 2, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "mno", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
	{"PRIMARY", "PRIMARY KEY", "mno_nn", "id", 1, sql.NullInt64{Int64: 0, Valid: false}, sql.NullString{String: "", Valid: false}, sql.NullString{String: "", Valid: false}},
}

var views = []View{
	{
		Table: "abc_v", ViewDefinition: "select `dbsql_test`.`abc`.`id` AS `id`,`dbsql_test`.`abc`.`code` AS `code`,`dbsql_test`.`abc`.`description` AS `description` from `dbsql_test`.`abc` order by `dbsql_test`.`abc`.`code`",
		CheckOption: "NONE", IsUpdatable: "YES", Definer: "testuser@localhost",
		SecurityType: "DEFINER", CharacterSetClient: "utf8", CollationConnection: "utf8_general_ci",
	},
	{
		Table: "defghi_v", ViewDefinition: "select `a`.`id` AS `aid`,`b`.`id` AS `bid`,`a`.`d_datetime` AS `d_datetime`,`a`.`size` AS `size`,`b`.`stuff` AS `stuff` from `dbsql_test`.`def` `a` join `dbsql_test`.`ghi` `b` where (`a`.`id` = `b`.`def_id`) order by `a`.`id`,`a`.`size`,`b`.`def_id`",
		CheckOption: "NONE", IsUpdatable: "YES", Definer: "testuser@localhost",
		SecurityType: "DEFINER", CharacterSetClient: "utf8", CollationConnection: "utf8_general_ci",
	},
}

func TestMain(m *testing.M) {
	db, err := New(server, user, password, testDB)
	if err != nil {
		panic(err)
	}
	//defer TeardownTestDB(db.(*DB)) // this always tries to run, that way a partial setup is still torndown
	err = SetupTestDB(db.(*DB))
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestTables(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.GetTables()
	if err != nil {
		t.Errorf("unexpected error getting table information: %s", err)
		return
	}
	tables := m.Tables()
	for i, v := range tables {
		tbl, ok := v.(*Table)
		if !ok {
			t.Errorf("%s: assertion error; was not a Table", tableDefs[i].name)
		}
		if tbl.Name() != tableDefs[i].name {
			t.Errorf("name: got %q want %q", tbl.name, tableDefs[i].name)
			continue
		}
		if tbl.r != tableDefs[i].r {
			t.Errorf("%s.r: got %q want %q", tbl.name, tbl.r, tableDefs[i].r)
			continue
		}
		if tbl.StructName() != tableDefs[i].structName {
			t.Errorf("%s.StructName: got %q want %q", tbl.name, tbl.StructName(), tableDefs[i].structName)
			continue
		}
		if tbl.schema != tableDefs[i].schema {
			t.Errorf("%s.Schema: got %q want %q", tbl.name, tbl.schema, tableDefs[i].schema)
			continue
		}
		if tbl.Typ != tableDefs[i].Typ {
			t.Errorf("%s.Type: got %q want %q", tbl.name, tbl.Typ, tableDefs[i].Typ)
			continue
		}
		if tbl.Engine.Valid != tableDefs[i].Engine.Valid {
			t.Errorf("%s.Engine.Valid: got %t want %t", tbl.name, tbl.Engine.Valid, tableDefs[i].Engine.Valid)
			continue
		}
		if tbl.Engine.Valid {
			if tbl.Engine.String != tableDefs[i].Engine.String {
				t.Errorf("%s.Engine.String: got %q want %q", tbl.name, tbl.Engine.String, tableDefs[i].Engine.String)
				continue
			}
		}
		if tbl.collation.Valid != tableDefs[i].collation.Valid {
			t.Errorf("%s.Collation.Valid: got %t, want %t", tbl.name, tbl.collation.Valid, tableDefs[i].collation.Valid)
			continue
		}
		if tbl.collation.Valid {
			if tbl.collation.String != tableDefs[i].collation.String {
				t.Errorf("%s.Collation.String: got %q want %q", tbl.name, tbl.collation.String, tableDefs[i].collation.String)
				continue
			}
		}
		if tbl.Comment != tableDefs[i].Comment {
			t.Errorf("%s.Comment: got %q want %q", tbl.name, tbl.Comment, tableDefs[i].Comment)
			continue
		}
		// handle columns
		for j, col := range tbl.columns {
			if col.Name != tableDefs[i].columns[j].Name {
				t.Errorf("%s:%s COLUMN_NAME: got %q want %q", tbl.name, col.Name, col.Name, tableDefs[i].columns[j].Name)
				continue
			}
			if col.OrdinalPosition != tableDefs[i].columns[j].OrdinalPosition {
				t.Errorf("%s.%s ORDINAL_POSITION: got %q want %q", tbl.name, col.Name, col.OrdinalPosition, tableDefs[i].columns[j].OrdinalPosition)
				continue
			}
			if col.Default.Valid != tableDefs[i].columns[j].Default.Valid {
				t.Errorf("%s.%s DEFAULT Valid: got %t want %t", tbl.name, col.Name, col.Default.Valid, tableDefs[i].columns[j].Default.Valid)
				continue
			}
			if col.Default.Valid {
				if col.Default.String != tableDefs[i].columns[j].Default.String {
					t.Errorf("%s.%s DEFAULT String: got %s want %s", tbl.name, col.Name, col.Default.String, tableDefs[i].columns[j].Default.String)
				}
				continue
			}
			if col.IsNullable != tableDefs[i].columns[j].IsNullable {
				t.Errorf("%s.%s IS_NULLABLE: got %q want %q", tbl.name, col.Name, col.IsNullable, tableDefs[i].columns[j].IsNullable)
				continue
			}
			if col.DataType != tableDefs[i].columns[j].DataType {
				t.Errorf("%s.%s DATA_TYPE: got %q want %q", tbl.name, col.Name, col.DataType, tableDefs[i].columns[j].DataType)
				continue
			}
			if col.CharMaxLen.Valid != tableDefs[i].columns[j].CharMaxLen.Valid {
				t.Errorf("%s.%s CHARACTER_MAXIMUM_LENGTH Valid: got %t want %t", tbl.name, col.Name, col.CharMaxLen.Valid, tableDefs[i].columns[j].CharMaxLen.Valid)
				continue
			}
			if col.CharMaxLen.Valid {
				if col.CharMaxLen.Int64 != tableDefs[i].columns[j].CharMaxLen.Int64 {
					t.Errorf("%s.%s CHARACTER_MAXIMUM_LENGTH Int64: got %v want %v", tbl.name, col.Name, col.CharMaxLen.Int64, tableDefs[i].columns[j].CharMaxLen.Int64)
				}
				continue
			}
			if col.CharOctetLen.Valid != tableDefs[i].columns[j].CharOctetLen.Valid {
				t.Errorf("%s.%s CHARACTER_OCTET_LENGTH Valid: got %t want %t", tbl.name, col.Name, col.CharOctetLen.Valid, tableDefs[i].columns[j].CharOctetLen.Valid)
				continue
			}
			if col.CharOctetLen.Valid {
				if col.CharOctetLen.Int64 != tableDefs[i].columns[j].CharOctetLen.Int64 {
					t.Errorf("%s.%s CHARACTER_OCTET_LENGTH Int64: got %v want %v", tbl.name, col.Name, col.CharOctetLen.Int64, tableDefs[i].columns[j].CharOctetLen.Int64)
				}
				continue
			}
			if col.NumericPrecision.Valid != tableDefs[i].columns[j].NumericPrecision.Valid {
				t.Errorf("%s.%s NUMERIC_PRECISION Valid: got %t want %t", tbl.name, col.Name, col.NumericPrecision.Valid, tableDefs[i].columns[j].NumericPrecision.Valid)
				continue
			}
			if col.NumericPrecision.Valid {
				if col.NumericPrecision.Int64 != tableDefs[i].columns[j].NumericPrecision.Int64 {
					t.Errorf("%s.%s NUMERIC_PRECISION Int64: got %v want %v", tbl.name, col.Name, col.NumericPrecision.Int64, tableDefs[i].columns[j].NumericPrecision.Int64)
				}
				continue
			}
			if col.NumericScale.Valid != tableDefs[i].columns[j].NumericScale.Valid {
				t.Errorf("%s.%s NUMERIC_SCALE Valid: got %t want %t", tbl.name, col.Name, col.NumericScale.Valid, tableDefs[i].columns[j].NumericScale.Valid)
				continue
			}
			if col.NumericScale.Valid {
				if col.NumericScale.Int64 == tableDefs[i].columns[j].NumericScale.Int64 {
					t.Errorf("%s.%s NUMERIC_SCALE Int64: got %v want %v", tbl.name, col.Name, col.NumericScale.Int64, tableDefs[i].columns[j].NumericScale.Int64)
				}
				continue
			}
			if col.CharacterSet.Valid != tableDefs[i].columns[j].CharacterSet.Valid {
				t.Errorf("%s.%s CHARACTER_SET_NAME Valid: got %t want %t", tbl.name, col.Name, col.CharacterSet.Valid, tableDefs[i].columns[j].CharacterSet.Valid)
				continue
			}
			if col.CharacterSet.Valid {
				if col.CharacterSet.String != tableDefs[i].columns[j].CharacterSet.String {
					t.Errorf("%s.%s CHARACTER_SET_NAME String: got %s want %s", tbl.name, col.Name, col.CharacterSet.String, tableDefs[i].columns[j].CharacterSet.String)
				}
				continue
			}
			if col.Collation.Valid != tableDefs[i].columns[j].Collation.Valid {
				t.Errorf("%s.%s COLLATION_NAME Valid: got %t want %t", tbl.name, col.Name, col.Collation.Valid, tableDefs[i].columns[j].Collation.Valid)
				continue
			}
			if col.Collation.Valid {
				if col.Collation.String == tableDefs[i].columns[j].Collation.String {
					t.Errorf("%s.%s COLLATION_NAME String: got %s want %s", tbl.name, col.Name, col.Collation.String, tableDefs[i].columns[j].Collation.String)
				}
				continue
			}
			if col.Typ != tableDefs[i].columns[j].Typ {
				t.Errorf("%s.%s COLUMN_TYPE: got %q want %q", tbl.name, col.Name, col.Typ, tableDefs[i].columns[j].Typ)
				continue
			}
			if col.Key != tableDefs[i].columns[j].Key {
				t.Errorf("%s.%s COLUMN_KEY: got %q want %q", tbl.name, col.Name, col.Key, tableDefs[i].columns[j].Key)
				continue
			}
			if col.Extra != tableDefs[i].columns[j].Extra {
				t.Errorf("%s.%s EXTRA: got %q want %q", tbl.name, col.Name, col.Extra, tableDefs[i].columns[j].Extra)
				continue
			}
			if col.Privileges != tableDefs[i].columns[j].Privileges {
				t.Errorf("%s.%s PRIVILEGES: got %q want %q", tbl.name, col.Name, col.Privileges, tableDefs[i].columns[j].Privileges)
				continue
			}
			if col.Comment != tableDefs[i].columns[j].Comment {
				t.Errorf("%s.%s COMMENT: got %q want %q", tbl.name, col.Name, col.Comment, tableDefs[i].columns[j].Comment)
				continue
			}
			if col.fieldName != tableDefs[i].columns[j].fieldName {
				t.Errorf("%s.%s fieldName: got %q want %q", tbl.name, col.Name, col.fieldName, tableDefs[i].columns[j].fieldName)
				continue
			}
		}
	}
}

func TestIndexes(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.GetIndexes()
	if err != nil {
		t.Errorf("unexpected error getting index information: %s", err)
		return
	}
	for i, ndx := range m.(*DB).indexes {
		if ndx.Table != indexes[i].Table {
			t.Errorf("%s.%s.%d.Table: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Table, indexes[i].Table)
			continue
		}
		if ndx.NonUnique != indexes[i].NonUnique {
			t.Errorf("%s.%s.%d.NonUnique: got %d want %d", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.NonUnique, indexes[i].NonUnique)
			continue
		}
		if ndx.Schema != indexes[i].Schema {
			t.Errorf("%s.%s.%d.Schema: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Schema, indexes[i].Schema)
			continue
		}
		if ndx.name != indexes[i].name {
			t.Errorf("%s.%s.%d.Name: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.name, indexes[i].name)
			continue
		}
		if ndx.SeqInIndex != indexes[i].SeqInIndex {
			t.Errorf("%s.%s.%d.SeqInIndex: got %d want %d", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.SeqInIndex, indexes[i].SeqInIndex)
			continue
		}
		if ndx.Column != indexes[i].Column {
			t.Errorf("%s.%s.%d.Column: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Column, indexes[i].Column)
			continue
		}
		if ndx.Collation.Valid != indexes[i].Collation.Valid {
			t.Errorf("%s.%s.%d.Collation.Valid: got %t want %t", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Collation.Valid, indexes[i].Collation.Valid)
			continue
		}
		if ndx.Collation.Valid {
			if ndx.Collation.String != indexes[i].Collation.String {
				t.Errorf("%s.%s.%d.Collation.String: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Collation.String, indexes[i].Collation.String)
				continue
			}
		}
		if ndx.Cardinality.Valid != indexes[i].Cardinality.Valid {
			t.Errorf("%s.%s.%d.Cardinality.Valid: got %t want %t", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Cardinality.Valid, indexes[i].Cardinality.Valid)
			continue
		}
		if ndx.Cardinality.Valid {
			if ndx.Cardinality.Int64 != indexes[i].Cardinality.Int64 {
				t.Errorf("%s.%s.%d.Cardinality.Int64: got %d want %d", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Cardinality.Int64, indexes[i].Cardinality.Int64)
				continue
			}
		}
		if ndx.SubPart.Valid != indexes[i].SubPart.Valid {
			t.Errorf("%s.%s.%d.SubPart.Valid: got %t want %t", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.SubPart.Valid, indexes[i].SubPart.Valid)
			continue
		}
		if ndx.SubPart.Valid {
			if ndx.SubPart.Int64 != indexes[i].SubPart.Int64 {
				t.Errorf("%s.%s.%d.SubPart.Int64: got %d want %d", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.SubPart.Int64, indexes[i].SubPart.Int64)
				continue
			}
		}
		if ndx.Packed.Valid != indexes[i].Packed.Valid {
			t.Errorf("%s.%s.%d.Packed.Valid: got %t want %t", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Packed.Valid, indexes[i].Packed.Valid)
			continue
		}
		if ndx.Packed.Valid {
			if ndx.Packed.String != indexes[i].Packed.String {
				t.Errorf("%s.%s.%d.Packed.String: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Packed.String, indexes[i].Packed.String)
				continue
			}
		}
		if ndx.Nullable != indexes[i].Nullable {
			t.Errorf("%s.%s.%d.Nullable: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Nullable, indexes[i].Nullable)
			continue
		}
		if ndx.Type != indexes[i].Type {
			t.Errorf("%s.%s.%d.Type: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Type, indexes[i].Type)
			continue
		}
		if ndx.Comment.Valid != indexes[i].Comment.Valid {
			t.Errorf("%s.%s.%d.Comment.Valid: got %t want %t", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Comment.Valid, indexes[i].Comment.Valid)
			continue
		}
		if ndx.Comment.Valid {
			if ndx.Packed.String != indexes[i].Packed.String {
				t.Errorf("%s.%s.%d.Comment.String: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.Comment.String, indexes[i].Comment.String)
				continue
			}
		}
		if ndx.IndexComment != indexes[i].IndexComment {
			t.Errorf("%s.%s.%d.IndexComment: got %s want %s", ndx.Table, ndx.name, ndx.SeqInIndex, ndx.IndexComment, indexes[i].IndexComment)
			continue
		}
	}
}

func TestGetConstraints(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.GetConstraints()
	if err != nil {
		t.Errorf("unexpected error getting key information: %s", err)
		return
	}
	// Check key info
	for i, k := range m.(*DB).constraints {
		if k.Name != constraints[i].Name {
			t.Errorf("%s.%d.Name: got %s want %s", constraints[i].Name, constraints[i].Seq, k.Name, constraints[i].Name)
			continue
		}
		if k.Type != constraints[i].Type {
			t.Errorf("%s.%d.Type: got %s want %s", constraints[i].Name, constraints[i].Seq, k.Type, constraints[i].Type)
			continue
		}
		if k.Table != constraints[i].Table {
			t.Errorf("%s.%d.Table: got %s want %s", constraints[i].Name, constraints[i].Seq, k.Table, constraints[i].Table)
			continue
		}
		if k.Column != constraints[i].Column {
			t.Errorf("%s.%d.Column: got %s want %s", constraints[i].Name, constraints[i].Seq, k.Column, constraints[i].Column)
			continue
		}
		if k.Seq != constraints[i].Seq {
			t.Errorf("%s.%d.Seq: got %d want %d", constraints[i].Name, constraints[i].Seq, k.Seq, constraints[i].Seq)
			continue
		}
		if k.USeq.Valid != constraints[i].USeq.Valid {
			t.Errorf("%s.%d.USeq.Valid: got %t want %t", constraints[i].Name, constraints[i].Seq, k.USeq.Valid, constraints[i].USeq.Valid)
			continue
		}
		if k.USeq.Valid {
			if k.USeq.Int64 != constraints[i].USeq.Int64 {
				t.Errorf("%s.%d.USeq.Int64: got %v want %v", constraints[i].Name, constraints[i].Seq, k.USeq.Int64, constraints[i].USeq.Int64)
				continue
			}
		}
		if k.RefTable.Valid != constraints[i].RefTable.Valid {
			t.Errorf("%s.%d.RefTable.Valid: got %t want %t", constraints[i].Name, constraints[i].Seq, k.RefTable.Valid, constraints[i].RefTable.Valid)
			continue
		}
		if k.RefTable.Valid {
			if k.RefTable.String != constraints[i].RefTable.String {
				t.Errorf("%s.%d.RefTable.String: got %s want %s", constraints[i].Name, constraints[i].Seq, k.RefTable.String, constraints[i].RefTable.String)
				continue
			}
		}
		if k.RefCol.Valid != constraints[i].RefCol.Valid {
			t.Errorf("%s.%d.RefCol.Valid: got %t want %t", constraints[i].Name, constraints[i].Seq, k.RefCol.Valid, constraints[i].RefCol.Valid)
			continue
		}
		if k.RefCol.Valid {
			if k.RefCol.String != constraints[i].RefCol.String {
				t.Errorf("%s.%d.RefCol.String: got %s want %s", constraints[i].Name, constraints[i].Seq, k.RefCol.String, constraints[i].RefCol.String)
				continue
			}
		}

	}
}

func TestViews(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.GetViews()
	if err != nil {
		t.Errorf("unexpected error getting index information: %s", err)
		return
	}
	vs := m.Views()
	for i, view := range vs {
		v := view.(*View)
		if v.Table != views[i].Table {
			t.Errorf("%s: got %s; want %s", views[i].Table, v.Table, views[i].Table)
			continue
		}
		if v.ViewDefinition != views[i].ViewDefinition {
			t.Errorf("%s.ViewDefinition: got %s; want %s", views[i].Table, v.ViewDefinition, views[i].ViewDefinition)
			continue
		}
		if v.CheckOption != views[i].CheckOption {
			t.Errorf("%s.CheckOption: got %s; want %s", views[i].Table, v.CheckOption, views[i].CheckOption)
			continue
		}
		if v.IsUpdatable != views[i].IsUpdatable {
			t.Errorf("%s.IsUpdatable: got %s; want %s", views[i].IsUpdatable, v.Table, views[i].IsUpdatable)
			continue
		}
		if v.Definer != views[i].Definer {
			t.Errorf("%s.Definer: got %s; want %s", views[i].Table, v.Definer, views[i].Definer)
			continue
		}
		if v.SecurityType != views[i].SecurityType {
			t.Errorf("%s.SecurityType: got %s; want %s", views[i].Table, v.SecurityType, views[i].SecurityType)
			continue
		}
		if v.CharacterSetClient != views[i].CharacterSetClient {
			t.Errorf("%s.CharacterSetClient: got %s; want %s", views[i].Table, v.CharacterSetClient, views[i].CharacterSetClient)
			continue
		}
		if v.CollationConnection != views[i].CollationConnection {
			t.Errorf("%s.CollationConnection: got %s; want %s", views[i].Table, v.CollationConnection, views[i].CollationConnection)
			continue
		}
	}
}

func TestColumns(t *testing.T) {
	expected := []struct {
		name string
		cols []string
	}{
		{name: "abc", cols: []string{"id", "code", "description", "tiny", "small", "medium", "ger", "big", "cost", "created"}},
		{name: "abc_nn", cols: []string{"id", "code", "description", "tiny", "small", "medium", "ger", "big", "cost", "created"}},
		{name: "abc_v", cols: []string{"id", "code", "description"}},
		{name: "def", cols: []string{"id", "d_date", "d_datetime", "d_time", "d_year", "size", "a_set"}},
		{name: "def_nn", cols: []string{"id", "d_date", "d_datetime", "d_time", "d_year", "size", "a_set"}},
		{name: "defghi_v", cols: []string{"aid", "bid", "d_datetime", "size", "stuff"}},
		{name: "ghi", cols: []string{"id", "val", "def_id", "def_datetime", "tiny_stuff", "stuff", "med_stuff", "long_stuff"}},
		{name: "ghi_nn", cols: []string{"id", "val", "def_id", "def_datetime", "tiny_stuff", "stuff", "med_stuff", "long_stuff"}},
		{name: "jkl", cols: []string{"id", "fid", "tiny_txt", "txt", "med_txt", "long_txt", "bin", "var_bin"}},
		{name: "jkl_nn", cols: []string{"id", "fid", "tiny_txt", "txt", "med_txt", "long_txt", "bin", "var_bin"}},
		{name: "mno", cols: []string{"id", "geo", "pt", "lstring", "poly", "multi_pt", "multi_lstring", "multi_polygon", "geo_collection"}},
		{name: "mno_nn", cols: []string{"id", "geo", "pt", "lstring", "poly", "multi_pt", "multi_lstring", "multi_polygon", "geo_collection"}},
	}
	for i, tbl := range tableDefs {
		cols := tbl.Columns()
		if !sliceEqual(cols, expected[i].cols) {
			t.Errorf("%s: got %v want %v", expected[i].name, cols, expected[i].cols)
		}
	}
}

func TestUpdateTables(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.Get()
	if err != nil {
		t.Errorf("unexpected error getting database information: %s", err)
		return
	}

	for i, tbl := range m.(*DB).tables {
		ndxs := tbl.Indexes()
		for j, ndx := range ndxs {
			if tableDefs[i].indexes[j].Type != ndx.Type {
				t.Errorf("Index: %d:%d: %s:%s.Type: got %v; want %v", i, j, tbl.Name(), ndx.Name, ndx.Type, tableDefs[i].indexes[j].Type)
				continue
			}
			if tableDefs[i].indexes[j].Primary != ndx.Primary {
				t.Errorf("Index: %d:%d: %s:%s.Primary: got %v; want %v", i, j, tbl.Name(), ndx.Name, ndx.Primary, tableDefs[i].indexes[j].Primary)
				continue
			}
			if tableDefs[i].indexes[j].Name != ndx.Name {
				t.Errorf("Index: %d:%d: %s:%s.Name: got %v; want %v", i, j, tbl.Name(), ndx.Name, ndx.Name, tableDefs[i].indexes[j].Name)
				continue
			}
			if tableDefs[i].indexes[j].Table != ndx.Table {
				t.Errorf("Index: %d:%d: %s:%s.Type: got %v; want %v", i, j, tbl.Name(), ndx.Name, ndx.Table, tableDefs[i].indexes[j].Table)
				continue
			}
			for k, col := range ndx.Cols {
				if tableDefs[i].indexes[j].Cols[k] != col {
					t.Errorf("Index: %d:%d:%d: %s:%s.Cols.%d: got %s; want %s", i, j, k, tbl.Name(), ndx.Name, k, col, tableDefs[i].indexes[j].Cols[k])
					continue
				}
			}
		}
		cons := tbl.Constraints()
		for j, c := range cons {
			if tableDefs[i].constraints[j].Type != c.Type {
				t.Errorf("Constraint: %d:%d: %s:%s.Type: got %v; want %v", i, j, tbl.Name(), c.Name, c.Type, tableDefs[i].constraints[j].Type)
				continue
			}
			if tableDefs[i].constraints[j].Name != c.Name {
				t.Errorf("Constraint: %d:%d: %s:%s.Name: got %v; want %v", i, j, tbl.Name(), c.Name, c.Name, tableDefs[i].constraints[j].Name)
				continue
			}
			if tableDefs[i].constraints[j].Table != c.Table {
				t.Errorf("Constraint: %d:%d: %s:%s.Table: got %v; want %v", i, j, tbl.Name(), c.Name, c.Table, tableDefs[i].constraints[j].Table)
				continue
			}
			for k, col := range c.Cols {
				if tableDefs[i].constraints[j].Cols[k] != col {
					t.Errorf("Constraint: %d:%d:%d: %s:%s.Cols.%d: got %s; want %s", i, j, k, tbl.Name(), c.Name, k, col, tableDefs[i].constraints[j].Cols[k])
					continue
				}
			}
			if c.RefTable != tableDefs[i].constraints[j].RefTable {
				t.Errorf("Constraint: %d:%d: %s:%s.RefTable: got %v; want %v", i, j, tbl.Name(), c.Name, c.RefTable, tableDefs[i].constraints[j].RefTable)
				continue
			}
			for k, col := range c.RefCols {
				if tableDefs[i].constraints[j].RefCols[k] != col {
					t.Errorf("Constraint: %d:%d:%d: %s:%s.RefCols.%d: got %s; want %s", i, j, k, tbl.Name(), c.Name, k, col, tableDefs[i].constraints[j].RefCols[k])
					continue
				}
			}
		}
	}
}

func TestSetReceiverName(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.Get()
	if err != nil {
		t.Errorf("unexpected error getting database information: %s", err)
		return
	}
	tbls := m.Tables()
	for i, v := range tbls {
		if v.(*Table).r != tableDefs[i].r {
			t.Errorf("%s: got %c want %c", tableDefs[i].name, v.(*Table).r, tableDefs[i].r)
		}
	}
}

func TestGenerateDefs(t *testing.T) {
	var buf bytes.Buffer
	for i, def := range tableDefs {
		buf.Reset()
		if i == 7 { // geospatial is not yet implemented; so skip
			break
		}
		err := def.Definition(&buf)
		if err != nil {
			t.Errorf("%s: %s", def.Name(), err)
		}
		if tableDefsString[i] != buf.String() {
			t.Errorf("%s: got %q; want %q", def.Name(), buf.String(), tableDefsString[i])
		}
	}
}

func TestStructDefs(t *testing.T) {
	var buf bytes.Buffer
	for i, def := range tableDefs {
		buf.Reset()
		if i == 7 { // geospatial is not yet implemented; so skip
			break
		}
		err := def.GoFmt(&buf)
		if err != nil {
			t.Errorf("%s: %s", def.Name(), err)
		}
		if structDefs[i] != buf.String() {
			t.Errorf("%s: got %q; want %q", def.Name(), buf.String(), structDefs[i])
		}
	}
}

func TestSelectSQLPK(t *testing.T) {
	expected := []string{
		"SELECT id, code, description, tiny, small, medium, ger, big, cost, created FROM abc WHERE id = ?",
		"SELECT id, code, description, tiny, small, medium, ger, big, cost, created FROM abc_nn WHERE id = ?",
		"", // views don't have a PK so nothing is generated
		"SELECT id, d_date, d_datetime, d_time, d_year, size, a_set FROM def WHERE id = ?",
		"SELECT id, d_date, d_datetime, d_time, d_year, size, a_set FROM def_nn WHERE id = ?",
		"", // views don't have a PK so nothing is generated
		"", // NO PK, no sql generated
		"", // NO PK, no sql generated
		"SELECT id, fid, tiny_txt, txt, med_txt, long_txt, bin, var_bin FROM jkl WHERE id = ? AND fid = ?",
		"SELECT id, fid, tiny_txt, txt, med_txt, long_txt, bin, var_bin FROM jkl_nn WHERE id = ? AND fid = ?",
		"SELECT id, geo, pt, lstring, poly, multi_pt, multi_lstring, multi_polygon, geo_collection FROM mno WHERE id = ?",
		"SELECT id, geo, pt, lstring, poly, multi_pt, multi_lstring, multi_polygon, geo_collection FROM mno_nn WHERE id = ?",
	}

	var buf bytes.Buffer
	for i, tbl := range tableDefs {
		buf.Reset()
		err := tbl.SelectSQLPK(&buf)
		if err != nil {
			t.Errorf("%d: unexpected error: got %q", i, err)
			continue
		}
		if buf.String() != expected[i] {
			t.Errorf("%d: got %q; want %q", i, buf.String(), expected[i])
		}
	}
}

func TestDeleteSQLPK(t *testing.T) {
	expected := []string{
		"DELETE FROM abc WHERE id = ?",
		"DELETE FROM abc_nn WHERE id = ?",
		"", // delete from views not supported
		"DELETE FROM def WHERE id = ?",
		"DELETE FROM def_nn WHERE id = ?",
		"", // delete from views not supported
		"", // NO PK, no sql generated
		"", // NO PK, no sql generated
		"DELETE FROM jkl WHERE id = ? AND fid = ?",
		"DELETE FROM jkl_nn WHERE id = ? AND fid = ?",
		"DELETE FROM mno WHERE id = ?",
		"DELETE FROM mno_nn WHERE id = ?",
	}
	var buf bytes.Buffer
	for i, tbl := range tableDefs {
		buf.Reset()
		err := tbl.DeleteSQLPK(&buf)
		if err != nil {
			t.Errorf("%d: unexpected error: got %q", i, err)
			continue
		}
		if buf.String() != expected[i] {
			t.Errorf("%d: got %q; want %q", i, buf.String(), expected[i])
		}
	}
}

func TestInsertSQL(t *testing.T) {
	expected := [][]byte{
		[]byte("INSERT INTO abc (id, code, description, tiny, small, medium, ger, big, cost, created) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO abc_nn (id, code, description, tiny, small, medium, ger, big, cost, created) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		nil, // INSERT views not supported.
		[]byte("INSERT INTO def (id, d_date, d_datetime, d_time, d_year, size, a_set) VALUES (?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO def_nn (id, d_date, d_datetime, d_time, d_year, size, a_set) VALUES (?, ?, ?, ?, ?, ?, ?)"),
		nil, // INSERT views not supported.
		[]byte("INSERT INTO ghi (id, val, def_id, def_datetime, tiny_stuff, stuff, med_stuff, long_stuff) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO ghi_nn (id, val, def_id, def_datetime, tiny_stuff, stuff, med_stuff, long_stuff) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO jkl (id, fid, tiny_txt, txt, med_txt, long_txt, bin, var_bin) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO jkl_nn (id, fid, tiny_txt, txt, med_txt, long_txt, bin, var_bin) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO mno (id, geo, pt, lstring, poly, multi_pt, multi_lstring, multi_polygon, geo_collection) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"),
		[]byte("INSERT INTO mno_nn (id, geo, pt, lstring, poly, multi_pt, multi_lstring, multi_polygon, geo_collection) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"),
	}

	for i, tbl := range tableDefs {
		sql, err := tbl.InsertSQL()
		if err != nil {
			t.Errorf("%d: unexpected error: got %q", i, err)
			continue
		}
		if bytes.Compare(sql, expected[i]) != 0 {
			t.Errorf("%d: got %q; want %q", i, sql, expected[i])
		}
	}
}

func SetupTestDB(m *DB) error {
	// Everything is ignored because we don't care if it exists. This is just in
	// case the it wasn't dropped in a prior test due to a panic or something.
	m.Conn.Exec(fmt.Sprintf("DROP DATABASE %s", testDB))

	_, err := m.Conn.Exec(fmt.Sprintf("CREATE DATABASE %s", testDB))
	if err != nil {
		return err
	}

	_, err = m.Conn.Exec(fmt.Sprintf("USE %s", testDB))
	if err != nil {
		return err
	}

	for _, v := range createTables {
		_, err := m.Conn.Exec(v)
		if err != nil {
			return err
		}
	}
	for _, v := range createViews {
		_, err := m.Conn.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func TeardownTestDB(m *DB) {
	// we don't care if this fails
	m.Conn.Exec(fmt.Sprintf("DROP DATABASE %s", testDB))
}

func sliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}
