package mysql

import (
	"database/sql"
	"fmt"
	"testing"
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
		a_set SET('a', 'b', 'c')
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE def_nn (
		id INT AUTO_INCREMENT PRIMARY KEY,
		d_date DATE NOT NULL,
		d_datetime DATETIME NOT NULL,
		d_time TIME NOT NULL,
		d_year YEAR NOT NULL,
		size ENUM('small', 'med', 'large') NOT NULL,
		a_set SET('1', '2', '3') NOT NULL
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE ghi (
		tiny_stuff TINYBLOB,
		stuff BLOB,
		med_stuff MEDIUMBLOB,
		long_stuff LONGBLOB
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE ghi_nn (
		tiny_stuff TINYBLOB NOT NULL,
		stuff BLOB NOT NULL,
		med_stuff MEDIUMBLOB NOT NULL,
		long_stuff LONGBLOB NOT NULL
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE jkl (
		id INT AUTO_INCREMENT PRIMARY KEY,
		tiny_txt TINYTEXT,
		txt TEXT,
		med_txt MEDIUMTEXT,
		long_txt LONGTEXT,
		bin BINARY(3),
		var_bin VARBINARY(12)
	)
	CHARACTER SET ascii COLLATE ascii_general_ci`,
	`CREATE TABLE jkl_nn (
		id INT AUTO_INCREMENT PRIMARY KEY,
		tiny_txt TINYTEXT NOT NULL,
		txt TEXT NOT NULL,
		med_txt MEDIUMTEXT NOT NULL,
		long_txt LONGTEXT NOT NULL,
		bin BINARY(3) NOT NULL,
		var_bin VARBINARY(12) NOT NULL
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

var tableDefs = []Table{
	Table{
		Name: "abc", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "code", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "char", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "char(12)",
				Key: "UNI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "description", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varchar", CharMaxLen: sql.NullInt64{Int64: 20, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 20, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "varchar(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny", OrdinalPosition: 4, Default: sql.NullString{String: "3", Valid: true},
				IsNullable: "YES", DataType: "tinyint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 3, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyint(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "small", OrdinalPosition: 5, Default: sql.NullString{String: "11", Valid: true},
				IsNullable: "YES", DataType: "smallint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 5, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "smallint(6)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "medium", OrdinalPosition: 6, Default: sql.NullString{String: "42", Valid: true},
				IsNullable: "YES", DataType: "mediumint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 7, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumint(9)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "ger", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "big", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "bigint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 19, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "bigint(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "cost", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "decimal", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "decimal(10,0)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "created", OrdinalPosition: 10, Default: sql.NullString{String: "CURRENT_TIMESTAMP", Valid: true},
				IsNullable: "NO", DataType: "timestamp", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "timestamp",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "latin1_swedish_ci", Comment: "",
	},
	Table{
		Name: "abc_nn", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "code", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "char", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "char(12)",
				Key: "UNI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "description", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varchar", CharMaxLen: sql.NullInt64{Int64: 20, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 20, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "latin1", Valid: true}, Collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Typ: "varchar(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinyint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 3, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyint(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "small", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "smallint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 5, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "smallint(6)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "medium", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 7, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumint(9)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "ger", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "big", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "bigint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 19, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "bigint(20)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "cost", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "decimal", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "decimal(10,0)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "created", OrdinalPosition: 10, Default: sql.NullString{String: "CURRENT_TIMESTAMP", Valid: true},
				IsNullable: "NO", DataType: "timestamp", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "timestamp",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "latin1_swedish_ci", Comment: "",
	},
	Table{
		Name: "def", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_date", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "date", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "date",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_datetime", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_time", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "time", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "time",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_year", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "year", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "year(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "size", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "enum", CharMaxLen: sql.NullInt64{Int64: 6, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 18, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "enum('small','medium','large')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "a_set", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "set", CharMaxLen: sql.NullInt64{Int64: 5, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 15, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "set('a','b','c')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "utf8_general_ci", Comment: "",
	},
	Table{
		Name: "def_nn", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_date", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "date", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "date",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_datetime", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_time", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "time", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "time",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "d_year", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "year", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "year(4)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "size", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "enum", CharMaxLen: sql.NullInt64{Int64: 5, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 15, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "enum('small','medium','large')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "a_set", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "set", CharMaxLen: sql.NullInt64{Int64: 5, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 15, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "set('a','b','c')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "utf8_general_ci", Comment: "",
	},
	Table{
		Name: "ghi", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "tiny_stuff", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "tinyblob", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "stuff", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_stuff", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "mediumblob", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_stuff", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "longblob", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "longblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "utf8_general_ci", Comment: "",
	},
	Table{
		Name: "ghi_nn", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "tiny_stuff", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinyblob", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "stuff", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_stuff", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumblob", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_stuff", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "longblob", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "longblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "utf8_general_ci", Comment: "",
	},
	Table{
		Name: "jkl", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny_txt", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "tinytext", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "tinytext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "txt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "text", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "text",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_txt", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "mediumtext", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "mediumtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_txt", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "longtext", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "longtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "bin", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "binary", CharMaxLen: sql.NullInt64{Int64: 3, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 3, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "binary(3)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "var_bin", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "varbinary", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "varbinary(12)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "ascii_general_ci", Comment: "",
	},
	Table{
		Name: "jkl_nn", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny_txt", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinytext", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "tinytext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "txt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "text", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "text",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_txt", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumtext", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "mediumtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_txt", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "longtext", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "longtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "bin", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "binary", CharMaxLen: sql.NullInt64{Int64: 3, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 3, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "binary(3)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "var_bin", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varbinary", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "varbinary(12)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "ascii_general_ci", Comment: "",
	},
	Table{
		Name: "mno", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "geo", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "geometry", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometry",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "pt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "point", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "point",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "lstring", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "linestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "linestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "poly", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "polygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "polygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "multi_pt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "multipoint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipoint",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "multi_lstring", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "multilinestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multilinestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "multi_polygon", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "multipolygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipolygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "geo_collection", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "geometrycollection", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometrycollection",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "utf8_general_ci", Comment: "",
	},
	Table{
		Name: "mno_nn", Schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "geo", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "geometry", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometry",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "pt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "point", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "point",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "lstring", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "linestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "linestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "poly", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "polygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "polygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "multi_pt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "multipoint", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipoint",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "multi_lstring", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "multilinestring", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multilinestring",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "multi_polygon", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "multipolygon", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "multipolygon",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "geo_collection", OrdinalPosition: 9, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "geometrycollection", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "geometrycollection",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: "InnoDB",
		Collation: "utf8_general_ci", Comment: "",
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
	`type Ghi struct {
	TinyStuff []byte
	Stuff []byte
	MedStuff []byte
	LongStuff []byte
}
`,
	`type GhiNn struct {
	TinyStuff []byte
	Stuff []byte
	MedStuff []byte
	LongStuff []byte
}
`,
	`type Jkl struct {
	ID int32
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
	TinyTxt []byte
	Txt []byte
	MedTxt []byte
	LongTxt []byte
	Bin []byte
	VarBin []byte
}
`,
}

func TestMain(m *testing.M) {
	db, err := NewMySQLDB(server, user, password, testDB)
	if err != nil {
		panic(err)
		return
	}
	defer TeardownTestDB(db) // this always tries to run, that way a partial setup is still torndown
	err = SetupTestDB(db)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestGetTables(t *testing.T) {
	m, err := NewMySQLDB(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	tables, err := m.GetTables()
	if err != nil {
		t.Errorf("unexpected error getting table information: %s", err)
		return
	}
	for i, v := range tables {
		if v.Name != tableDefs[i].Name {
			t.Errorf("name: got %q want %q", v.Name, tableDefs[i].Name)
			continue
		}
		if v.Typ != tableDefs[i].Typ {
			t.Errorf("Type: got %q want %q", v.Typ, tableDefs[i].Typ)
			continue
		}
		if v.Engine != tableDefs[i].Engine {
			t.Errorf("Engine: got %q want %q", v.Engine, tableDefs[i].Engine)
			continue
		}
		if v.Collation != tableDefs[i].Collation {
			t.Errorf("Collation: got %q want %q", v.Collation, tableDefs[i].Collation)
			continue
		}
		if v.Comment != tableDefs[i].Comment {
			t.Errorf("Comment: got %q want %q", v.Comment, tableDefs[i].Comment)
			continue
		}
		// handle columns
		for j, col := range v.Columns {
			if col.Name != tableDefs[i].Columns[j].Name {
				t.Errorf("%s:%s COLUMN_NAME: got %q want %q", v.Name, col.Name, col.Name, tableDefs[i].Columns[j].Name)
				continue
			}
			if col.OrdinalPosition != tableDefs[i].Columns[j].OrdinalPosition {
				t.Errorf("%s.%s ORDINAL_POSITION: got %q want %q", v.Name, col.Name, col.OrdinalPosition, tableDefs[i].Columns[j].OrdinalPosition)
				continue
			}
			if col.Default.Valid != tableDefs[i].Columns[j].Default.Valid {
				t.Errorf("%s.%s DEFAULT Valid: got %t want %t", v.Name, col.Name, col.Default.Valid, tableDefs[i].Columns[j].Default.Valid)
				continue
			}
			if col.Default.Valid {
				if col.Default.String != tableDefs[i].Columns[j].Default.String {
					t.Errorf("%s.%s DEFAULT String: got %s want %s", v.Name, col.Name, col.Default.String, tableDefs[i].Columns[j].Default.String)
				}
				continue
			}
			if col.IsNullable != tableDefs[i].Columns[j].IsNullable {
				t.Errorf("%s.%s IS_NULLABLE: got %q want %q", v.Name, col.Name, col.IsNullable, tableDefs[i].Columns[j].IsNullable)
				continue
			}
			if col.DataType != tableDefs[i].Columns[j].DataType {
				t.Errorf("%s.%s DATA_TYPE: got %q want %q", v.Name, col.Name, col.DataType, tableDefs[i].Columns[j].DataType)
				continue
			}
			if col.CharMaxLen.Valid != tableDefs[i].Columns[j].CharMaxLen.Valid {
				t.Errorf("%s.%s CHARACTER_MAXIMUM_LENGTH Valid: got %t want %t", v.Name, col.Name, col.CharMaxLen.Valid, tableDefs[i].Columns[j].CharMaxLen.Valid)
				continue
			}
			if col.CharMaxLen.Valid {
				if col.CharMaxLen.Int64 != tableDefs[i].Columns[j].CharMaxLen.Int64 {
					t.Errorf("%s.%s CHARACTER_MAXIMUM_LENGTH Int64: got %v want %v", v.Name, col.Name, col.CharMaxLen.Int64, tableDefs[i].Columns[j].CharMaxLen.Int64)
				}
				continue
			}
			if col.CharOctetLen.Valid != tableDefs[i].Columns[j].CharOctetLen.Valid {
				t.Errorf("%s.%s CHARACTER_OCTET_LENGTH Valid: got %t want %t", v.Name, col.Name, col.CharOctetLen.Valid, tableDefs[i].Columns[j].CharOctetLen.Valid)
				continue
			}
			if col.CharOctetLen.Valid {
				if col.CharOctetLen.Int64 != tableDefs[i].Columns[j].CharOctetLen.Int64 {
					t.Errorf("%s.%s CHARACTER_OCTET_LENGTH Int64: got %v want %v", v.Name, col.Name, col.CharOctetLen.Int64, tableDefs[i].Columns[j].CharOctetLen.Int64)
				}
				continue
			}
			if col.NumericPrecision.Valid != tableDefs[i].Columns[j].NumericPrecision.Valid {
				t.Errorf("%s.%s NUMERIC_PRECISION Valid: got %t want %t", v.Name, col.Name, col.NumericPrecision.Valid, tableDefs[i].Columns[j].NumericPrecision.Valid)
				continue
			}
			if col.NumericPrecision.Valid {
				if col.NumericPrecision.Int64 != tableDefs[i].Columns[j].NumericPrecision.Int64 {
					t.Errorf("%s.%s NUMERIC_PRECISION Int64: got %v want %v", v.Name, col.Name, col.NumericPrecision.Int64, tableDefs[i].Columns[j].NumericPrecision.Int64)
				}
				continue
			}
			if col.NumericScale.Valid != tableDefs[i].Columns[j].NumericScale.Valid {
				t.Errorf("%s.%s NUMERIC_SCALE Valid: got %t want %t", v.Name, col.Name, col.NumericScale.Valid, tableDefs[i].Columns[j].NumericScale.Valid)
				continue
			}
			if col.NumericScale.Valid {
				if col.NumericScale.Int64 == tableDefs[i].Columns[j].NumericScale.Int64 {
					t.Errorf("%s.%s NUMERIC_SCALE Int64: got %v want %v", v.Name, col.NumericScale.Int64, tableDefs[i].Columns[j].NumericScale.Int64)
				}
				continue
			}
			if col.CharacterSet.Valid != tableDefs[i].Columns[j].CharacterSet.Valid {
				t.Errorf("%s.%s CHARACTER_SET_NAME Valid: got %t want %t", v.Name, col.Name, col.CharacterSet.Valid, tableDefs[i].Columns[j].CharacterSet.Valid)
				continue
			}
			if col.CharacterSet.Valid {
				if col.CharacterSet.String != tableDefs[i].Columns[j].CharacterSet.String {
					t.Errorf("%s.%s CHARACTER_SET_NAME String: got %s want %s", v.Name, col.Name, col.CharacterSet.String, tableDefs[i].Columns[j].CharacterSet.String)
				}
				continue
			}
			if col.Collation.Valid != tableDefs[i].Columns[j].Collation.Valid {
				t.Errorf("%s.%s COLLATION_NAME Valid: got %t want %t", v.Name, col.Name, col.Collation.Valid, tableDefs[i].Columns[j].Collation.Valid)
				continue
			}
			if col.Collation.Valid {
				if col.Collation.String == tableDefs[i].Columns[j].Collation.String {
					t.Errorf("%s.%s COLLATION_NAME String: got %s want %s", v.Name, col.Name, col.Collation.String, tableDefs[i].Columns[j].Collation.String)
				}
				continue
			}
			if col.Typ != tableDefs[i].Columns[j].Typ {
				t.Errorf("%s.%s COLUMN_TYPE: got %q want %q", v.Name, col.Name, col.Typ, tableDefs[i].Columns[j].Typ)
				continue
			}
			if col.Key != tableDefs[i].Columns[j].Key {
				t.Errorf("%s.%s COLUMN_KEY: got %q want %q", v.Name, col.Name, col.Key, tableDefs[i].Columns[j].Key)
				continue
			}
			if col.Extra != tableDefs[i].Columns[j].Extra {
				t.Errorf("%s.%s EXTRA: got %q want %q", v.Name, col.Name, col.Extra, tableDefs[i].Columns[j].Extra)
				continue
			}
			if col.Privileges != tableDefs[i].Columns[j].Privileges {
				t.Errorf("%s.%s PRIVILEGES: got %q want %q", v.Name, col.Name, col.Privileges, tableDefs[i].Columns[j].Privileges)
				continue
			}
			if col.Comment != tableDefs[i].Columns[j].Comment {
				t.Errorf("%s.%s COMMENT: got %q want %q", v.Name, col.Name, col.Comment, tableDefs[i].Columns[j].Comment)
				continue
			}
		}
	}
}

func TestGenerateDefs(t *testing.T) {
	for i, def := range tableDefs {
		if i == 7 { // geospatial is not yet implemented; so skip
			break
		}
		d, err := def.Go()
		if err != nil {
			t.Error("%s: %s", def.Name, err)
		}
		if tableDefsString[i] != string(d) {
			t.Errorf("%s: got %q; want %q", def.Name, string(d), tableDefsString[i])
		}
	}
}

func SetupTestDB(m *MySQLDB) error {
	_, err := m.DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", testDB))
	//_, err := m.DB.Exec("CREATE DATABASE IF NOT EXISTS ?", m.dbName)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(fmt.Sprintf("USE %s", testDB))
	if err != nil {
		return err
	}

	for _, v := range createTables {
		_, err := m.DB.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func TeardownTestDB(m *MySQLDB) {
	_, err := m.DB.Exec(fmt.Sprintf("DROP DATABASE %s", testDB))
	if err != nil {
		panic(err)
	}
}
