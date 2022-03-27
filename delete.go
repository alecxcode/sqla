package sqla

import (
	"database/sql"
	"log"
	"sort"
)

// VerifyRemovalPermissions makes query to check if a user has the right to delete an object.
// The returned result will be truth if either an Owner matches id in column 'Creator' or have AdminPrivileges is true.
// RemoveAllowed flag defines if any remove allowed at all by non-admin user.
// This function is somewhat specific to EDM project. You might need to modify it for your app.
func VerifyRemovalPermissions(db *sql.DB, DBType byte, table string, Owner int, AdminPrivileges bool, RemoveAllowed bool, ids []int) bool {
	if AdminPrivileges {
		return true
	}
	if !RemoveAllowed {
		return false
	}
	var argsCounter int
	var args, argstoAppend []interface{}
	var sqlids []int

	argsCounter++
	var sq = "SELECT ID FROM " + table + " WHERE Creator = " + MakeParam(DBType, argsCounter) + " "
	args = append(args, Owner)

	argsCounter, sq, argstoAppend = BuildSQLIN(DBType, sq, argsCounter, "ID", ids)
	args = append(args, argstoAppend...)

	if DEBUG {
		log.Println(sq, args)
	}
	rows, err := db.Query(sq, args...)
	if err != nil {
		log.Println(currentFunction()+":", err)
	}
	defer rows.Close()
	var ID sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&ID)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		sqlids = append(sqlids, int(ID.Int64))
	}
	sort.Ints(ids)
	sort.Ints(sqlids)
	if intSlicesEqual(ids, sqlids) {
		return true
	}

	return false
}

// DeleteObject just deletes one specific object which id is in column specified.
func DeleteObject(db *sql.DB, DBType byte, table string, column string, id int) (rowsaff int) {
	var sq = "DELETE FROM " + table + " WHERE " + column + " = " + MakeParam(DBType, 1)
	if DEBUG {
		log.Println(sq, id)
	}
	res, err := db.Exec(sq, id)
	if err != nil {
		log.Println(currentFunction()+":", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		log.Println(currentFunction()+":", err)
	}
	rowsaff = int(ra)
	return rowsaff
}

// DeleteObjects just deletes any object which id is present in ids list and in column specified.
func DeleteObjects(db *sql.DB, DBType byte, table string, column string, ids []int) (rowsaff int) {

	var sq = "DELETE FROM " + table + " "
	var args, argstoAppend []interface{}
	var argsCounter int
	var res sql.Result
	var err error

	if len(ids) > 0 {
		argsCounter, sq, argstoAppend = BuildSQLIN(DBType, sq, argsCounter, column, ids)
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
