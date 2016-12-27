package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Abc struct {
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

var abcSelect = `SELECT id, code, description
	, tiny, small, medium
	, ger, big, cost
	, created
FROM abc
ORDER BY id`

type AbcNn struct {
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

var abcNnSelect = `SELECT id, code, description
	, tiny, small, medium
	, ger, big, cost
	, created
FROM abc_nn
ORDER BY id`

type Def struct {
	ID        int32
	DDate     mysql.NullTime
	DDatetime mysql.NullTime
	DTime     sql.NullString
	DYear     sql.NullString
	Size      []byte
	ASet      []byte
}

var defSelect = `SELECT id, d_date, d_datetime, d_time, d_year, size, a_set
FROM def
ORDER BY id`

type DefNn struct {
	ID        int32
	DDate     mysql.NullTime
	DDatetime mysql.NullTime
	DTime     string
	DYear     string
	Size      string
	ASet      string
}

var defNnSelect = `SELECT id, d_date, d_datetime, d_time, d_year, size, a_set
FROM def_nn
ORDER BY id`

type Ghi struct {
	TinyStuff []byte
	Stuff     []byte
	MedStuff  []byte
	LongStuff []byte
}

var ghiSelect = `SELECT tiny_stuff, stuff, med_stuff, long_stuff FROM ghi`

type GhiNn struct {
	TinyStuff []byte
	Stuff     []byte
	MedStuff  []byte
	LongStuff []byte
}

var ghiNnSelect = `SELECT tiny_stuff, stuff, med_stuff, long_stuff FROM ghi_nn`

type Jkl struct {
	ID      int32
	TinyTxt []byte
	Txt     []byte
	MedTxt  []byte
	LongTxt []byte
	Bin     []byte
	VarBin  []byte
}

var jklSelect = `SELECT ID, tiny_txt, txt, med_txt, long_txt, bin, var_bin
FROM jkl
ORDER BY id`

type JklNn struct {
	ID      int32
	TinyTxt []byte
	Txt     []byte
	MedTxt  []byte
	LongTxt []byte
	Bin     []byte
	VarBin  []byte
}

var jklNnSelect = `SELECT ID, tiny_txt, txt, med_txt, long_txt, bin, var_bin
FROM jkl_nn
ORDER BY id`

func NewMySQLDB(server, user, password, database string) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, database))
}

func main() {
	conn, err := NewMySQLDB("locahlhost", "testuser", "testuser", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	abcs, err := selectABC(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range abcs {
		fmt.Printf("%#v\n", v)
	}

	abcNns, err := selectABCNn(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range abcNns {
		fmt.Printf("%#v\n", v)
	}

	defs, err := selectDEF(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range defs {
		fmt.Printf("%#v\n", v)
	}

	defNns, err := selectDEFNn(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range defNns {
		fmt.Printf("%#v\n", v)
	}

	ghis, err := selectGHI(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range ghis {
		fmt.Printf("%#v\n", v)
	}

	ghiNns, err := selectGHINn(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range ghiNns {
		fmt.Printf("%#v\n", v)
	}

	jkls, err := selectJKL(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range jkls {
		fmt.Printf("%#v\n", v)
	}

	jklNns, err := selectJKLNn(conn)
	if err != nil {
		log.Print(err)
		return
	}
	for _, v := range jklNns {
		fmt.Printf("%#v\n", v)
	}
}

func selectABC(conn *sql.DB) ([]Abc, error) {
	var abcs []Abc
	rows, err := conn.Query(abcSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var abc Abc
		err = rows.Scan(&abc.ID, &abc.Code, &abc.Description,
			&abc.Tiny, &abc.Small, &abc.Medium,
			&abc.Ger, &abc.Big, &abc.Cost,
			&abc.Created)
		if err != nil {
			return nil, err
		}
		abcs = append(abcs, abc)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return abcs, nil
}

func selectABCNn(conn *sql.DB) ([]AbcNn, error) {
	var abcs []AbcNn
	rows, err := conn.Query(abcNnSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var abc AbcNn
		err = rows.Scan(&abc.ID, &abc.Code, &abc.Description,
			&abc.Tiny, &abc.Small, &abc.Medium,
			&abc.Ger, &abc.Big, &abc.Cost,
			&abc.Created)
		if err != nil {
			return nil, err
		}
		abcs = append(abcs, abc)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return abcs, nil
}

func selectDEF(conn *sql.DB) ([]Def, error) {
	var defs []Def
	rows, err := conn.Query(defSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var def Def
		err = rows.Scan(&def.ID, &def.DDate, &def.DDatetime, &def.DTime, &def.DYear, &def.Size, &def.ASet)
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return defs, nil
}

func selectDEFNn(conn *sql.DB) ([]DefNn, error) {
	var defs []DefNn
	rows, err := conn.Query(defNnSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var def DefNn
		err = rows.Scan(&def.ID, &def.DDate, &def.DDatetime, &def.DTime, &def.DYear, &def.Size, &def.ASet)
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return defs, nil
}

func selectGHI(conn *sql.DB) ([]Ghi, error) {
	var ghis []Ghi
	rows, err := conn.Query(ghiSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ghi Ghi
		err = rows.Scan(&ghi.TinyStuff, &ghi.Stuff, &ghi.MedStuff, &ghi.LongStuff)
		if err != nil {
			return nil, err
		}
		ghis = append(ghis, ghi)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ghis, nil
}

func selectGHINn(conn *sql.DB) ([]GhiNn, error) {
	var ghis []GhiNn
	rows, err := conn.Query(ghiNnSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ghi GhiNn
		err = rows.Scan(&ghi.TinyStuff, &ghi.Stuff, &ghi.MedStuff, &ghi.LongStuff)
		if err != nil {
			return nil, err
		}
		ghis = append(ghis, ghi)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ghis, nil
}

func selectJKL(conn *sql.DB) ([]Jkl, error) {
	var jkls []Jkl
	rows, err := conn.Query(jklSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var jkl Jkl
		err = rows.Scan(&jkl.ID, &jkl.TinyTxt, &jkl.Txt, &jkl.MedTxt, &jkl.LongTxt, &jkl.Bin, &jkl.VarBin)
		if err != nil {
			return nil, err
		}
		jkls = append(jkls, jkl)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return jkls, nil
}

func selectJKLNn(conn *sql.DB) ([]JklNn, error) {
	var jkls []JklNn
	rows, err := conn.Query(jklNnSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var jkl JklNn
		err = rows.Scan(&jkl.ID, &jkl.TinyTxt, &jkl.Txt, &jkl.MedTxt, &jkl.LongTxt, &jkl.Bin, &jkl.VarBin)
		if err != nil {
			return nil, err
		}
		jkls = append(jkls, jkl)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return jkls, nil
}
