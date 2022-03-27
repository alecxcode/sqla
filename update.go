package sqla

import (
	"database/sql"
	"log"
)

// UpdateObject creates an SQL statement and executes it to update an object in the specified table.
// The function returns the number of affected rows. Update will be done on the object where column 'ID' contains ID value.
func UpdateObject(db *sql.DB, DBType byte, table string, iargs []anyT, ID int) (rowsaff int) {

	const (
		I = 0
		B = 1
		F = 2
		S = 3
		N = 4
	)

	var colvalpairs string
	var counter int
	var args []interface{}
	for j := 0; j < len(iargs); j++ {
		counter++
		if j > 0 {
			colvalpairs += ", "
		}
		colvalpairs += iargs[j].c + " = " + MakeParam(DBType, counter)
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
	counter++
	args = append(args, ID)
	sq := "UPDATE " + table + " SET " + colvalpairs + " WHERE ID = " + MakeParam(DBType, counter)

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
	rowsaff = int(ra)

	return rowsaff

}

// UpdateMultipleWithOneInt updates with val the column of an object which id is present in ids list and in 'ID' column. Rows which already have val in the column will not be updated. If necessary you can provide timestamp and a column for timestamp; if you don't need to update any timestamp column use empty string as the argument for that column.
func UpdateMultipleWithOneInt(db *sql.DB, DBType byte, table string, column string, val int, timecol string, timestamp int64, ids []int) (rowsaff int) {
	var sq = "UPDATE " + table + " SET " + column + " = " + MakeParam(DBType, 1) + " "
	var args, argstoAppend []interface{}
	args = append(args, val)
	var argsCounter = 1
	var res sql.Result
	var err error

	if timecol != "" {
		argsCounter++
		sq += ", " + timecol + " = " + MakeParam(DBType, argsCounter) + " "
		args = append(args, timestamp)
	}

	argsCounter++
	sq += "WHERE " + column + " <> " + MakeParam(DBType, argsCounter) + " "
	args = append(args, val)

	if len(ids) > 0 {
		argsCounter, sq, argstoAppend = BuildSQLIN(DBType, sq, argsCounter, "ID", ids)
		args = append(args, argstoAppend...)
		if DEBUG {
			log.Println(sq, args)
		}
		res, err = db.Exec(sq, args...)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		rowsaff = int(ra)
	}
	return rowsaff
}

// SetToNull set to NULL the column of any object which has a value from a list in that column.
func SetToNull(db *sql.DB, DBType byte, table string, column string, list []int) (rowsaff int) {

	var sq = "UPDATE " + table + " SET " + column + " = NULL "
	var args, argstoAppend []interface{}
	var argsCounter int
	var res sql.Result
	var err error

	if len(list) > 0 {
		argsCounter, sq, argstoAppend = BuildSQLIN(DBType, sq, argsCounter, column, list)
		args = append(args, argstoAppend...)
		if DEBUG {
			log.Println(sq, args)
		}
		res, err = db.Exec(sq, args...)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		rowsaff = int(ra)
	}
	return rowsaff
}

// UpdateSingleInt creates an SQL statement and executes it to update only one value of an object in database.
// It executes UpdateObject.
func UpdateSingleInt(db *sql.DB, DBType byte, table string, column string, valueToSet int, ID int) (rowsaff int) {
	var args AnyTslice
	args = args.AppendInt(column, valueToSet)
	rowsaff = UpdateObject(db, DBType, table, args, ID)
	return rowsaff
}

// UpdateSingleStr creates an SQL statement and executes it to update only one value of an object in database.
// It executes UpdateObject.
func UpdateSingleStr(db *sql.DB, DBType byte, table string, column string, valueToSet string, ID int) (rowsaff int) {
	var args AnyTslice
	args = args.AppendStringOrNil(column, valueToSet)
	rowsaff = UpdateObject(db, DBType, table, args, ID)
	return rowsaff
}

// UpdateSingleJSONStruct creates an SQL statement and executes it to update only one value of an object in database.
// It executes UpdateObject.
func UpdateSingleJSONStruct(db *sql.DB, DBType byte, table string, column string, valueToSet interface{}, ID int) (rowsaff int) {
	var args AnyTslice
	args = args.AppendJSONStruct(column, valueToSet)
	rowsaff = UpdateObject(db, DBType, table, args, ID)
	return rowsaff
}

// UpdateSingleJSONListStr creates an SQL statement and executes it to update only one value of an object in database.
// It executes UpdateObject.
func UpdateSingleJSONListStr(db *sql.DB, DBType byte, table string, column string, valueToSet []string, ID int) (rowsaff int) {
	var args AnyTslice
	args = args.AppendJSONList(column, valueToSet)
	rowsaff = UpdateObject(db, DBType, table, args, ID)
	return rowsaff
}

// UpdateSingleJSONListInt creates an SQL statement and executes it to update only one value of an object in database.
// It executes UpdateObject.
func UpdateSingleJSONListInt(db *sql.DB, DBType byte, table string, column string, valueToSet []int, ID int) (rowsaff int) {
	var args AnyTslice
	args = args.AppendJSONListInt(column, valueToSet)
	rowsaff = UpdateObject(db, DBType, table, args, ID)
	return rowsaff
}
