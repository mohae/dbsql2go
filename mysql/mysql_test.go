package mysql

import (
	"database/sql"
	"fmt"
	"sort"
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
	Table{
		name: "abc", schema: "dbsql_test",
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
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "abc_nn", schema: "dbsql_test",
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
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "latin1_swedish_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "abc_v", schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "0", Valid: true},
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
		},
		Typ: "VIEW", Engine: sql.NullString{String: "", Valid: false},
		collation: sql.NullString{String: "", Valid: false}, Comment: "VIEW",
	},
	Table{
		name: "def", schema: "dbsql_test",
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
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "def_nn", schema: "dbsql_test",
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
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "defghi_v", schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "aid", OrdinalPosition: 1, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "bid", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "PRI", Extra: "", Privileges: "select,insert,update,references",
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
				Name: "size", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "enum", CharMaxLen: sql.NullInt64{Int64: 6, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 18, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "utf8", Valid: true}, Collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Typ: "enum('small','medium','large')",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "stuff", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "VIEW", Engine: sql.NullString{String: "", Valid: false},
		collation: sql.NullString{String: "", Valid: false}, Comment: "VIEW",
	},
	Table{
		name: "ghi", schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "val", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "def_id", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "def_datetime", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny_stuff", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "tinyblob", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "stuff", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_stuff", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "mediumblob", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_stuff", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "longblob", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "longblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "ghi_nn", schema: "dbsql_test",
		Columns: []Column{
			Column{
				Name: "id", OrdinalPosition: 1, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "val", OrdinalPosition: 2, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "def_id", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "def_datetime", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "datetime", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "datetime",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny_stuff", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinyblob", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "tinyblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "stuff", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "blob", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "blob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_stuff", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumblob", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "mediumblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_stuff", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "longblob", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "longblob",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "jkl", schema: "dbsql_test",
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
				Name: "fid", OrdinalPosition: 2, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "YES", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny_txt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "tinytext", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "tinytext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "txt", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "text", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "text",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_txt", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "mediumtext", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "mediumtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_txt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "longtext", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "longtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "bin", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "binary", CharMaxLen: sql.NullInt64{Int64: 3, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 3, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "binary(3)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "var_bin", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "YES", DataType: "varbinary", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "varbinary(12)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "jkl_nn", schema: "dbsql_test",
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
				Name: "fid", OrdinalPosition: 2, Default: sql.NullString{String: "0", Valid: true},
				IsNullable: "NO", DataType: "int", CharMaxLen: sql.NullInt64{Int64: 0, Valid: false},
				CharOctetLen: sql.NullInt64{Int64: 0, Valid: false}, NumericPrecision: sql.NullInt64{Int64: 10, Valid: true}, NumericScale: sql.NullInt64{Int64: 0, Valid: true},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "int(11)",
				Key: "MUL", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "tiny_txt", OrdinalPosition: 3, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "tinytext", CharMaxLen: sql.NullInt64{Int64: 255, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 255, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "tinytext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "txt", OrdinalPosition: 4, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "text", CharMaxLen: sql.NullInt64{Int64: 65535, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 65535, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "text",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "med_txt", OrdinalPosition: 5, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "mediumtext", CharMaxLen: sql.NullInt64{Int64: 16777215, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 16777215, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "mediumtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "long_txt", OrdinalPosition: 6, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "longtext", CharMaxLen: sql.NullInt64{Int64: 4294967295, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 4294967295, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "ascii", Valid: true}, Collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Typ: "longtext",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "bin", OrdinalPosition: 7, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "binary", CharMaxLen: sql.NullInt64{Int64: 3, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 3, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "binary(3)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
			Column{
				Name: "var_bin", OrdinalPosition: 8, Default: sql.NullString{String: "", Valid: false},
				IsNullable: "NO", DataType: "varbinary", CharMaxLen: sql.NullInt64{Int64: 12, Valid: true},
				CharOctetLen: sql.NullInt64{Int64: 12, Valid: true}, NumericPrecision: sql.NullInt64{Int64: 0, Valid: false}, NumericScale: sql.NullInt64{Int64: 0, Valid: false},
				CharacterSet: sql.NullString{String: "", Valid: false}, Collation: sql.NullString{String: "", Valid: false}, Typ: "varbinary(12)",
				Key: "", Extra: "", Privileges: "select,insert,update,references",
				Comment: "",
			},
		},
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "ascii_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "mno", schema: "dbsql_test",
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
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
	},
	Table{
		name: "mno_nn", schema: "dbsql_test",
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
		Typ: "BASE TABLE", Engine: sql.NullString{String: "InnoDB", Valid: true},
		collation: sql.NullString{String: "utf8_general_ci", Valid: true}, Comment: "",
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

var fmtdTableDefsString = []string{
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
`,
}

var indexes = []Index{
	{
		TableName: "abc", NonUnique: 0, Schema: "dbsql_test", name: "code",
		SeqInIndex: 1, ColumnName: "code", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "abc", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "abc_nn", NonUnique: 0, Schema: "dbsql_test", name: "code",
		SeqInIndex: 1, ColumnName: "code", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "abc_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "def", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "def", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 2, ColumnName: "d_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "def", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "def_nn", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "def_nn", NonUnique: 1, Schema: "dbsql_test", name: "id",
		SeqInIndex: 2, ColumnName: "d_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "def_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "ghi", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 1, ColumnName: "def_id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "ghi", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 2, ColumnName: "def_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "ghi", NonUnique: 1, Schema: "dbsql_test", name: "val",
		SeqInIndex: 1, ColumnName: "val", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "YES", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "ghi_nn", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 1, ColumnName: "def_id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "ghi_nn", NonUnique: 1, Schema: "dbsql_test", name: "fk_def",
		SeqInIndex: 2, ColumnName: "def_datetime", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "ghi_nn", NonUnique: 1, Schema: "dbsql_test", name: "val",
		SeqInIndex: 1, ColumnName: "val", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "jkl", NonUnique: 1, Schema: "dbsql_test", name: "fid",
		SeqInIndex: 1, ColumnName: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "jkl", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "jkl", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 2, ColumnName: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "jkl_nn", NonUnique: 1, Schema: "dbsql_test", name: "fid",
		SeqInIndex: 1, ColumnName: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "jkl_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "jkl_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 2, ColumnName: "fid", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "mno", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
	{
		TableName: "mno_nn", NonUnique: 0, Schema: "dbsql_test", name: "PRIMARY",
		SeqInIndex: 1, ColumnName: "id", Collation: sql.NullString{String: "A", Valid: true}, Cardinality: sql.NullInt64{Int64: 0, Valid: true},
		SubPart: sql.NullInt64{Int64: 0, Valid: false}, Packed: sql.NullString{String: "", Valid: false}, Nullable: "", Type: "BTREE",
		Comment: sql.NullString{String: "", Valid: true}, IndexComment: "",
	},
}

var keys = []Key{
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
		TableName: "abc_v", ViewDefinition: "select `dbsql_test`.`abc`.`id` AS `id`,`dbsql_test`.`abc`.`code` AS `code`,`dbsql_test`.`abc`.`description` AS `description` from `dbsql_test`.`abc` order by `dbsql_test`.`abc`.`code`",
		CheckOption: "NONE", IsUpdatable: "YES", Definer: "testuser@localhost",
		SecurityType: "DEFINER", CharacterSetClient: "utf8", CollationConnection: "utf8_general_ci",
	},
	{
		TableName: "defghi_v", ViewDefinition: "select `a`.`id` AS `aid`,`b`.`id` AS `bid`,`a`.`d_datetime` AS `d_datetime`,`a`.`size` AS `size`,`b`.`stuff` AS `stuff` from `dbsql_test`.`def` `a` join `dbsql_test`.`ghi` `b` where (`a`.`id` = `b`.`def_id`) order by `a`.`id`,`a`.`size`,`b`.`def_id`",
		CheckOption: "NONE", IsUpdatable: "YES", Definer: "testuser@localhost",
		SecurityType: "DEFINER", CharacterSetClient: "utf8", CollationConnection: "utf8_general_ci",
	},
}

func TestMain(m *testing.M) {
	db, err := New(server, user, password, testDB)
	if err != nil {
		panic(err)
		return
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
		for j, col := range tbl.Columns {
			if col.Name != tableDefs[i].Columns[j].Name {
				t.Errorf("%s:%s COLUMN_NAME: got %q want %q", tbl.name, col.Name, col.Name, tableDefs[i].Columns[j].Name)
				continue
			}
			if col.OrdinalPosition != tableDefs[i].Columns[j].OrdinalPosition {
				t.Errorf("%s.%s ORDINAL_POSITION: got %q want %q", tbl.name, col.Name, col.OrdinalPosition, tableDefs[i].Columns[j].OrdinalPosition)
				continue
			}
			if col.Default.Valid != tableDefs[i].Columns[j].Default.Valid {
				t.Errorf("%s.%s DEFAULT Valid: got %t want %t", tbl.name, col.Name, col.Default.Valid, tableDefs[i].Columns[j].Default.Valid)
				continue
			}
			if col.Default.Valid {
				if col.Default.String != tableDefs[i].Columns[j].Default.String {
					t.Errorf("%s.%s DEFAULT String: got %s want %s", tbl.name, col.Name, col.Default.String, tableDefs[i].Columns[j].Default.String)
				}
				continue
			}
			if col.IsNullable != tableDefs[i].Columns[j].IsNullable {
				t.Errorf("%s.%s IS_NULLABLE: got %q want %q", tbl.name, col.Name, col.IsNullable, tableDefs[i].Columns[j].IsNullable)
				continue
			}
			if col.DataType != tableDefs[i].Columns[j].DataType {
				t.Errorf("%s.%s DATA_TYPE: got %q want %q", tbl.name, col.Name, col.DataType, tableDefs[i].Columns[j].DataType)
				continue
			}
			if col.CharMaxLen.Valid != tableDefs[i].Columns[j].CharMaxLen.Valid {
				t.Errorf("%s.%s CHARACTER_MAXIMUM_LENGTH Valid: got %t want %t", tbl.name, col.Name, col.CharMaxLen.Valid, tableDefs[i].Columns[j].CharMaxLen.Valid)
				continue
			}
			if col.CharMaxLen.Valid {
				if col.CharMaxLen.Int64 != tableDefs[i].Columns[j].CharMaxLen.Int64 {
					t.Errorf("%s.%s CHARACTER_MAXIMUM_LENGTH Int64: got %v want %v", tbl.name, col.Name, col.CharMaxLen.Int64, tableDefs[i].Columns[j].CharMaxLen.Int64)
				}
				continue
			}
			if col.CharOctetLen.Valid != tableDefs[i].Columns[j].CharOctetLen.Valid {
				t.Errorf("%s.%s CHARACTER_OCTET_LENGTH Valid: got %t want %t", tbl.name, col.Name, col.CharOctetLen.Valid, tableDefs[i].Columns[j].CharOctetLen.Valid)
				continue
			}
			if col.CharOctetLen.Valid {
				if col.CharOctetLen.Int64 != tableDefs[i].Columns[j].CharOctetLen.Int64 {
					t.Errorf("%s.%s CHARACTER_OCTET_LENGTH Int64: got %v want %v", tbl.name, col.Name, col.CharOctetLen.Int64, tableDefs[i].Columns[j].CharOctetLen.Int64)
				}
				continue
			}
			if col.NumericPrecision.Valid != tableDefs[i].Columns[j].NumericPrecision.Valid {
				t.Errorf("%s.%s NUMERIC_PRECISION Valid: got %t want %t", tbl.name, col.Name, col.NumericPrecision.Valid, tableDefs[i].Columns[j].NumericPrecision.Valid)
				continue
			}
			if col.NumericPrecision.Valid {
				if col.NumericPrecision.Int64 != tableDefs[i].Columns[j].NumericPrecision.Int64 {
					t.Errorf("%s.%s NUMERIC_PRECISION Int64: got %v want %v", tbl.name, col.Name, col.NumericPrecision.Int64, tableDefs[i].Columns[j].NumericPrecision.Int64)
				}
				continue
			}
			if col.NumericScale.Valid != tableDefs[i].Columns[j].NumericScale.Valid {
				t.Errorf("%s.%s NUMERIC_SCALE Valid: got %t want %t", tbl.name, col.Name, col.NumericScale.Valid, tableDefs[i].Columns[j].NumericScale.Valid)
				continue
			}
			if col.NumericScale.Valid {
				if col.NumericScale.Int64 == tableDefs[i].Columns[j].NumericScale.Int64 {
					t.Errorf("%s.%s NUMERIC_SCALE Int64: got %v want %v", tbl.name, col.Name, col.NumericScale.Int64, tableDefs[i].Columns[j].NumericScale.Int64)
				}
				continue
			}
			if col.CharacterSet.Valid != tableDefs[i].Columns[j].CharacterSet.Valid {
				t.Errorf("%s.%s CHARACTER_SET_NAME Valid: got %t want %t", tbl.name, col.Name, col.CharacterSet.Valid, tableDefs[i].Columns[j].CharacterSet.Valid)
				continue
			}
			if col.CharacterSet.Valid {
				if col.CharacterSet.String != tableDefs[i].Columns[j].CharacterSet.String {
					t.Errorf("%s.%s CHARACTER_SET_NAME String: got %s want %s", tbl.name, col.Name, col.CharacterSet.String, tableDefs[i].Columns[j].CharacterSet.String)
				}
				continue
			}
			if col.Collation.Valid != tableDefs[i].Columns[j].Collation.Valid {
				t.Errorf("%s.%s COLLATION_NAME Valid: got %t want %t", tbl.name, col.Name, col.Collation.Valid, tableDefs[i].Columns[j].Collation.Valid)
				continue
			}
			if col.Collation.Valid {
				if col.Collation.String == tableDefs[i].Columns[j].Collation.String {
					t.Errorf("%s.%s COLLATION_NAME String: got %s want %s", tbl.name, col.Name, col.Collation.String, tableDefs[i].Columns[j].Collation.String)
				}
				continue
			}
			if col.Typ != tableDefs[i].Columns[j].Typ {
				t.Errorf("%s.%s COLUMN_TYPE: got %q want %q", tbl.name, col.Name, col.Typ, tableDefs[i].Columns[j].Typ)
				continue
			}
			if col.Key != tableDefs[i].Columns[j].Key {
				t.Errorf("%s.%s COLUMN_KEY: got %q want %q", tbl.name, col.Name, col.Key, tableDefs[i].Columns[j].Key)
				continue
			}
			if col.Extra != tableDefs[i].Columns[j].Extra {
				t.Errorf("%s.%s EXTRA: got %q want %q", tbl.name, col.Name, col.Extra, tableDefs[i].Columns[j].Extra)
				continue
			}
			if col.Privileges != tableDefs[i].Columns[j].Privileges {
				t.Errorf("%s.%s PRIVILEGES: got %q want %q", tbl.name, col.Name, col.Privileges, tableDefs[i].Columns[j].Privileges)
				continue
			}
			if col.Comment != tableDefs[i].Columns[j].Comment {
				t.Errorf("%s.%s COMMENT: got %q want %q", tbl.name, col.Name, col.Comment, tableDefs[i].Columns[j].Comment)
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
		if ndx.TableName != indexes[i].TableName {
			t.Errorf("%s.%s.%d.Tablename: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.TableName, indexes[i].TableName)
			continue
		}
		if ndx.NonUnique != indexes[i].NonUnique {
			t.Errorf("%s.%s.%d.NonUnique: got %d want %d", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.NonUnique, indexes[i].NonUnique)
			continue
		}
		if ndx.Schema != indexes[i].Schema {
			t.Errorf("%s.%s.%d.Schema: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Schema, indexes[i].Schema)
			continue
		}
		if ndx.name != indexes[i].name {
			t.Errorf("%s.%s.%d.Name: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Name, indexes[i].Name)
			continue
		}
		if ndx.SeqInIndex != indexes[i].SeqInIndex {
			t.Errorf("%s.%s.%d.SeqInIndex: got %d want %d", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.SeqInIndex, indexes[i].SeqInIndex)
			continue
		}
		if ndx.ColumnName != indexes[i].ColumnName {
			t.Errorf("%s.%s.%d.ColumnName: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.ColumnName, indexes[i].ColumnName)
			continue
		}
		if ndx.Collation.Valid != indexes[i].Collation.Valid {
			t.Errorf("%s.%s.%d.Collation.Valid: got %t want %t", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Collation.Valid, indexes[i].Collation.Valid)
			continue
		}
		if ndx.Collation.Valid {
			if ndx.Collation.String != indexes[i].Collation.String {
				t.Errorf("%s.%s.%d.Collation.String: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Collation.String, indexes[i].Collation.String)
				continue
			}
		}
		if ndx.Cardinality.Valid != indexes[i].Cardinality.Valid {
			t.Errorf("%s.%s.%d.Cardinality.Valid: got %t want %t", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Cardinality.Valid, indexes[i].Cardinality.Valid)
			continue
		}
		if ndx.Cardinality.Valid {
			if ndx.Cardinality.Int64 != indexes[i].Cardinality.Int64 {
				t.Errorf("%s.%s.%d.Cardinality.Int64: got %d want %d", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Cardinality.Int64, indexes[i].Cardinality.Int64)
				continue
			}
		}
		if ndx.SubPart.Valid != indexes[i].SubPart.Valid {
			t.Errorf("%s.%s.%d.SubPart.Valid: got %t want %t", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.SubPart.Valid, indexes[i].SubPart.Valid)
			continue
		}
		if ndx.SubPart.Valid {
			if ndx.SubPart.Int64 != indexes[i].SubPart.Int64 {
				t.Errorf("%s.%s.%d.SubPart.Int64: got %d want %d", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.SubPart.Int64, indexes[i].SubPart.Int64)
				continue
			}
		}
		if ndx.Packed.Valid != indexes[i].Packed.Valid {
			t.Errorf("%s.%s.%d.Packed.Valid: got %t want %t", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Packed.Valid, indexes[i].Packed.Valid)
			continue
		}
		if ndx.Packed.Valid {
			if ndx.Packed.String != indexes[i].Packed.String {
				t.Errorf("%s.%s.%d.Packed.String: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Packed.String, indexes[i].Packed.String)
				continue
			}
		}
		if ndx.Nullable != indexes[i].Nullable {
			t.Errorf("%s.%s.%d.Nullable: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Nullable, indexes[i].Nullable)
			continue
		}
		if ndx.Type != indexes[i].Type {
			t.Errorf("%s.%s.%d.Type: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Type, indexes[i].Type)
			continue
		}
		if ndx.Comment.Valid != indexes[i].Comment.Valid {
			t.Errorf("%s.%s.%d.Comment.Valid: got %t want %t", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Comment.Valid, indexes[i].Comment.Valid)
			continue
		}
		if ndx.Comment.Valid {
			if ndx.Packed.String != indexes[i].Packed.String {
				t.Errorf("%s.%s.%d.Comment.String: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.Comment.String, indexes[i].Comment.String)
				continue
			}
		}
		if ndx.IndexComment != indexes[i].IndexComment {
			t.Errorf("%s.%s.%d.IndexComment: got %s want %s", ndx.TableName, ndx.name, ndx.SeqInIndex, ndx.IndexComment, indexes[i].IndexComment)
			continue
		}
	}
}

func TestGetKeys(t *testing.T) {
	m, err := New(server, user, password, testDB)
	if err != nil {
		t.Errorf("unexpected connection error: %s", err)
		return
	}
	err = m.GetKeys()
	if err != nil {
		t.Errorf("unexpected error getting key information: %s", err)
		return
	}
	// Check key info
	for i, k := range m.(*DB).keys {
		if k.Name != keys[i].Name {
			t.Errorf("%s.%d.Name: got %s want %s", keys[i].Name, keys[i].Seq, k.Name, keys[i].Name)
			continue
		}
		if k.Type != keys[i].Type {
			t.Errorf("%s.%d.Type: got %s want %s", keys[i].Name, keys[i].Seq, k.Type, keys[i].Type)
			continue
		}
		if k.Table != keys[i].Table {
			t.Errorf("%s.%d.Table: got %s want %s", keys[i].Name, keys[i].Seq, k.Table, keys[i].Table)
			continue
		}
		if k.Column != keys[i].Column {
			t.Errorf("%s.%d.Column: got %s want %s", keys[i].Name, keys[i].Seq, k.Column, keys[i].Column)
			continue
		}
		if k.Seq != keys[i].Seq {
			t.Errorf("%s.%d.Seq: got %d want %d", keys[i].Name, keys[i].Seq, k.Seq, keys[i].Seq)
			continue
		}
		if k.USeq.Valid != keys[i].USeq.Valid {
			t.Errorf("%s.%d.USeq.Valid: got %t want %t", keys[i].Name, keys[i].Seq, k.USeq.Valid, keys[i].USeq.Valid)
			continue
		}
		if k.USeq.Valid {
			if k.USeq.Int64 != keys[i].USeq.Int64 {
				t.Errorf("%s.%d.USeq.Int64: got %s want %s", keys[i].Name, keys[i].Seq, k.USeq.Int64, keys[i].USeq.Int64)
				continue
			}
		}
		if k.RefTable.Valid != keys[i].RefTable.Valid {
			t.Errorf("%s.%d.RefTable.Valid: got %t want %t", keys[i].Name, keys[i].Seq, k.RefTable.Valid, keys[i].RefTable.Valid)
			continue
		}
		if k.RefTable.Valid {
			if k.RefTable.String != keys[i].RefTable.String {
				t.Errorf("%s.%d.RefTable.String: got %s want %s", keys[i].Name, keys[i].Seq, k.RefTable.String, keys[i].RefTable.String)
				continue
			}
		}
		if k.RefCol.Valid != keys[i].RefCol.Valid {
			t.Errorf("%s.%d.RefCol.Valid: got %t want %t", keys[i].Name, keys[i].Seq, k.RefCol.Valid, keys[i].RefCol.Valid)
			continue
		}
		if k.RefCol.Valid {
			if k.RefCol.String != keys[i].RefCol.String {
				t.Errorf("%s.%d.RefCol.String: got %s want %s", keys[i].Name, keys[i].Seq, k.RefCol.String, keys[i].RefCol.String)
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
		if v.TableName != views[i].TableName {
			t.Errorf("%s: got %s; want %s", views[i].TableName, v.TableName, views[i].TableName)
			continue
		}
		if v.ViewDefinition != views[i].ViewDefinition {
			t.Errorf("%s.ViewDefinition: got %s; want %s", views[i].TableName, v.ViewDefinition, views[i].ViewDefinition)
			continue
		}
		if v.CheckOption != views[i].CheckOption {
			t.Errorf("%s.CheckOption: got %s; want %s", views[i].TableName, v.CheckOption, views[i].CheckOption)
			continue
		}
		if v.IsUpdatable != views[i].IsUpdatable {
			t.Errorf("%s.IsUpdatable: got %s; want %s", views[i].IsUpdatable, v.TableName, views[i].IsUpdatable)
			continue
		}
		if v.Definer != views[i].Definer {
			t.Errorf("%s.Definer: got %s; want %s", views[i].TableName, v.Definer, views[i].Definer)
			continue
		}
		if v.SecurityType != views[i].SecurityType {
			t.Errorf("%s.SecurityType: got %s; want %s", views[i].TableName, v.SecurityType, views[i].SecurityType)
			continue
		}
		if v.CharacterSetClient != views[i].CharacterSetClient {
			t.Errorf("%s.CharacterSetClient: got %s; want %s", views[i].TableName, v.CharacterSetClient, views[i].CharacterSetClient)
			continue
		}
		if v.CollationConnection != views[i].CollationConnection {
			t.Errorf("%s.CollationConnection: got %s; want %s", views[i].TableName, v.CollationConnection, views[i].CollationConnection)
			continue
		}
	}
}

func TestColumnNames(t *testing.T) {
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
		cols := tbl.ColumnNames()
		if !sliceEqual(cols, expected[i].cols) {
			t.Errorf("%s: got %v want %v", expected[i].name, cols, expected[i].cols)
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
			t.Errorf("%s: got %q; want %q", def.Name(), string(d), tableDefsString[i])
		}
	}
}

func TestGenerateFmtdDefs(t *testing.T) {
	for i, def := range tableDefs {
		if i == 7 { // geospatial is not yet implemented; so skip
			break
		}
		d, err := def.GoFmt()
		if err != nil {
			t.Error("%s: %s", def.Name, err)
		}
		if fmtdTableDefsString[i] != string(d) {
			t.Errorf("%s: got %q; want %q", def.Name(), string(d), fmtdTableDefsString[i])
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
