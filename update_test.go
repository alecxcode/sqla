package sqla

import (
	"database/sql"
	"testing"
)

func TestUpdateObject(t *testing.T) {
	const (
		I = 0
		B = 1
		F = 2
		S = 3
	)
	var ID = 99
	var testobj = []anyT{
		{c: "RegNo", t: S, s: "#890-x"},
		{c: "RegDate", t: S, s: "2022-01-01"},
		{c: "DocType", t: I, i: 1},
	}
	var DBType byte
	DBType = SQLITE
	var db *sql.DB
	db = OpenSQLConnection(DBType, "file::memory:?cache=shared&_foreign_keys=true")
	db.Exec("CREATE TABLE documents (ID INTEGER PRIMARY KEY, RegNo TEXT, RegDate INTEGER, DocType INTEGER);")
	UpdateObject(db, DBType, "documents", testobj, ID)
	db.Close()
}
