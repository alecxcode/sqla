package sqla

import (
	"strconv"
	"strings"
)

// MakeParam takes database type (see constants) and parameter positional number to create ordial (or positional) placeholder for a parameter in your SQL statement.
// E.g. for the parameter #1 they are: $1 for SQLITE, @p1 for MSSQL, ? for MYSQL, :1 for ORACLE, $1 for POSTGRESQL. The function return only this placeholder.
func MakeParam(DBType byte, argsCounter int) string {
	switch DBType {
	case SQLITE:
		return "$" + strconv.Itoa(argsCounter)
	case MSSQL:
		return "@p" + strconv.Itoa(argsCounter)
	case MYSQL:
		return "?"
	case ORACLE:
		return ":" + strconv.Itoa(argsCounter)
	case POSTGRESQL:
		return "$" + strconv.Itoa(argsCounter)
	}
	return "?"
}

// BuildSQLIN makes and adds a part of SQL statement with all numbered parameters supplied in a valueList.
// sq argument should be complete SQL statement, as BuildSQLIN returns augmented statement, not part of it. argsCounter is required to define from what number to start count positional parameters.
// The added part of a statement will look something like 'WHERE/AND column IN(a list of placeholders by MakeParam func, e.g. $1, $2)'.
// BuildSQLIN returns counter as the number of added parameters for use in other routine and args as []interface{} of payload arguments which may be supplied to Go sql functions.
func BuildSQLIN(DBType byte, sq string, argsCounter int, column string, valueList []int) (counter int, resquery string, args []interface{}) {
	if strings.Contains(sq, "WHERE") {
		sq += "AND "
	} else {
		sq += "WHERE "
	}
	sq += column + " IN ("
	for i := 0; i < len(valueList); i++ {
		argsCounter++
		if i == 0 {
			sq += MakeParam(DBType, argsCounter)
		} else {
			sq += ", " + MakeParam(DBType, argsCounter)
		}
		args = append(args, valueList[i])
	}
	sq += ") "
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

// BuildSQLINNOT makes and adds a part of SQL statement with all numbered placeholders for arguments supplied in a valueList.
// sq argument should be complete SQL statement, as BuildSQLINNOT returns augmented statement, not part of it. argsCounter is required to define from what number to start count positional parameters.
// The added part of a statement will look something like 'WHERE/AND column NOT IN(a list of placeholders by MakeParam func, e.g. $1, $2)'.
// BuildSQLINNOT returns counter as the number of added parameters for use in other routine, statement, and args as []interface{} of payload arguments which may be supplied to Go sql functions.
func BuildSQLINNOT(DBType byte, sq string, argsCounter int, column string, valueList []int) (counter int, resquery string, args []interface{}) {
	if strings.Contains(sq, "WHERE") {
		sq += "AND "
	} else {
		sq += "WHERE "
	}
	sq += column + " NOT IN ("
	for i := 0; i < len(valueList); i++ {
		argsCounter++
		if i == 0 {
			sq += MakeParam(DBType, argsCounter)
		} else {
			sq += ", " + MakeParam(DBType, argsCounter)
		}
		args = append(args, valueList[i])
	}
	sq += ") "
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

// BuildSQLINOR makes and adds a part of SQL statement with all numbered placeholders for arguments supplied in a valueList.
// sq argument should be complete SQL statement, as BuildSQLINOR returns augmented statement, not part of it. argsCounter is required to define from what number to start count positional parameters.
// The added part of a statement will look something like 'WHERE/AND (column IN($1, $2) OR another_column IN($1, $2), etc...)'.
// Although the 'WHERE/AND (' part will be added only when FirstIter argument is true and the closing ')' will be added only when LastIter argument is true.
// BuildSQLINOR returns counter as the number of added parameters for use in other routine, statement, and args as []interface{} of payload arguments which may be supplied to Go sql functions.
func BuildSQLINOR(DBType byte, sq string, argsCounter int, column string, valueList []int, FirstIter bool, LastIter bool) (counter int, resquery string, args []interface{}) {
	if FirstIter {
		if strings.Contains(sq, "WHERE") {
			sq += "AND ("
		} else {
			sq += "WHERE ("
		}
	}
	if !FirstIter {
		sq += "OR "
	}
	sq += column + " IN ("
	for i := 0; i < len(valueList); i++ {
		argsCounter++
		if i == 0 {
			sq += MakeParam(DBType, argsCounter)
		} else {
			sq += ", " + MakeParam(DBType, argsCounter)
		}
		args = append(args, valueList[i])
	}
	sq += ") "
	if LastIter {
		sq += ") "
	}
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}
