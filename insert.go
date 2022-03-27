package sqla

import (
	"database/sql"
	"log"
)

// InsertObject creates an SQL statement and executes it to insert an object into the specified table.
// It returns the ID of created record and the number of affected rows. ID column should be named 'ID'.
func InsertObject(db *sql.DB, DBType byte, table string, iargs []anyT) (lastid int, rowsaff int) {

	const (
		I = 0
		B = 1
		F = 2
		S = 3
		N = 4
	)

	var columns string
	var values string
	var counter int
	var args []interface{}
	for j := 0; j < len(iargs); j++ {
		counter++
		if j > 0 {
			columns += ", "
			values += ", "
		}
		columns += iargs[j].c
		values += MakeParam(DBType, counter)
		switch iargs[j].t {
		case I:
			args = append(args, iargs[j].i)
		case B:
			args = append(args, iargs[j].b)
		case F:
			args = append(args, iargs[j].f)
		case S:
			args = append(args, iargs[j].s)
		case N:
			args = append(args, nil)
		}
	}

	if DBType == SQLITE || DBType == MYSQL {
		sq := "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ")"
		if DEBUG {
			log.Println(sq, args)
		}
		res, err := db.Exec(sq, args...)
		if err != nil {
			log.Println(currentFunction()+":", err)
			return
		}
		li, err := res.LastInsertId()
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		lastid = int(li)
		rowsaff = int(ra)
	}

	if DBType == MSSQL {
		sq := `DECLARE @virttable TABLE (NewID INTEGER);
		INSERT INTO ` + table + ` (` + columns + `) OUTPUT INSERTED.ID INTO @virttable VALUES (` + values + `);
		SELECT NewID FROM @virttable`
		if DEBUG {
			log.Println(sq, args)
		}
		row := db.QueryRow(sq, args...)
		var ra int
		var ID sql.NullInt64
		err := row.Scan(&ID)
		if err != nil {
			log.Println(currentFunction()+":", err)
			return
		}
		if ID.Valid == true {
			lastid = int(ID.Int64)
			ra++
		}
		rowsaff = int(ra)
	}

	if DBType == ORACLE {
		var li int64
		counter++
		args = append(args, &li)
		sq := "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ") RETURNING ID INTO " + MakeParam(DBType, counter)
		if DEBUG {
			log.Println(sq, args)
		}
		res, err := db.Exec(sq, args...)
		if err != nil {
			log.Println(currentFunction()+":", err)
			return
		}
		ra, err := res.RowsAffected()
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		lastid = int(li)
		rowsaff = int(ra)
	}

	if DBType == POSTGRESQL {
		sq := "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ") RETURNING ID"
		if DEBUG {
			log.Println(sq, args)
		}
		row := db.QueryRow(sq, args...)
		var ra int
		var ID sql.NullInt64
		err := row.Scan(&ID)
		if err != nil {
			log.Println(currentFunction()+":", err)
			return
		}
		if ID.Valid == true {
			lastid = int(ID.Int64)
			ra++
		}
		rowsaff = int(ra)
	}

	return lastid, rowsaff

}
