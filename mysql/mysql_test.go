package mysql

import (
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
	`CREATE TABLE def (
		id INT AUTO_INCREMENT PRIMARY KEY,
		d_date DATE,
		d_datetime DATETIME,
		d_time TIME,
		d_YEAR YEAR,
		size ENUM('small', 'medium', 'large'),
		aset SET('a', 'b', 'c')
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE ghi (
		tiny_stuff TINYBLOB,
		stuff BLOB,
		med_stuff MEDIUMBLOB,
		long_stuff LONGBLOB
	)
	CHARACTER SET utf8 COLLATE utf8_general_ci`,
	`CREATE TABLE jkl (
		id INT AUTO_INCREMENT PRIMARY KEY,
		tiny_txt TINYTEXT,
		txt TEXT,
		med_txt MEDIUMTEXT,
		long_txt LONGTEXT,
		bin BINARY(3),
		varbin VARBINARY(12)
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
	//	`CREATE TABLE pqr (
	//		id INT AUTO_INCREMENT PRIMARY KEY,
	//		jsn JSON DEFAULT NULL
	//	)`,
}

var expectedTables = []Table{
	Table{"abc", nil, "BASE TABLE", "InnoDB", "latin1_swedish_ci", ""},
	Table{"def", nil, "BASE TABLE", "InnoDB", "utf8_general_ci", ""},
	Table{"ghi", nil, "BASE TABLE", "InnoDB", "utf8_general_ci", ""},
	Table{"jkl", nil, "BASE TABLE", "InnoDB", "ascii_general_ci", ""},
	Table{"mno", nil, "BASE TABLE", "InnoDB", "utf8_general_ci", ""},
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
		if v.name != expectedTables[i].name {
			t.Errorf("name: got %q want %q", v.name, expectedTables[i].name)
			continue
		}
		if v.Type != expectedTables[i].Type {
			t.Errorf("Type: got %q want %q", v.Type, expectedTables[i].Type)
			continue
		}
		if v.Engine != expectedTables[i].Engine {
			t.Errorf("Engine: got %q want %q", v.Engine, expectedTables[i].Engine)
			continue
		}
		if v.Collation != expectedTables[i].Collation {
			t.Errorf("Collation: got %q want %q", v.Collation, expectedTables[i].Collation)
			continue
		}
		if v.Comment != expectedTables[i].Comment {
			t.Errorf("Comment: got %q want %q", v.Comment, expectedTables[i].Comment)
			continue
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
			fmt.Println(v)
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
