package sqla

import (
	"database/sql"
	"testing"
)

func TestInsertObject(t *testing.T) {
	const (
		I = 0
		B = 1
		F = 2
		S = 3
	)
	var testobj = []anyT{
		{c: "ID", t: I, i: 1},
		{c: "RegNo", t: S, s: "#890-x"},
		{c: "RegDate", t: S, s: "2022-01-01"},
	}
	var DBType byte
	DBType = SQLITE
	var db *sql.DB
	db = OpenSQLConnection(DBType, "file::memory:?cache=shared&_foreign_keys=true")
	db.Exec("CREATE TABLE documents (ID INTEGER PRIMARY KEY, RegNo TEXT, RegDate INTEGER, DocType INTEGER);")
	InsertObject(db, DBType, "documents", testobj)
	db.Close()
}
