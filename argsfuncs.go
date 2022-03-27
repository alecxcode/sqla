package sqla

import (
	"strconv"
	"strings"
)

func getRelationFromString(inprel string) (operator string) {
	switch inprel {
	case "eq", "=":
		operator = " = "
	case "gt", ">":
		operator = " > "
	case "lt", "<":
		operator = " < "
	case "gteq", ">=":
		operator = " >= "
	case "lteq", "<=":
		operator = " <= "
	case "noteq", "<>", "!=":
		operator = " <> "
	default:
		operator = " = "
	}
	return operator
}

func buildSQLCOMPARE(DBType byte, sq string, argsCounter int, column string, operator string, val interface{}) (counter int, resquery string, args []interface{}) {
	if strings.Contains(sq, "WHERE") {
		sq += "AND "
	} else {
		sq += "WHERE "
	}
	argsCounter++
	if operator == " <> " {
		sq += "(" + column + " IS NULL OR " + column + operator + MakeParam(DBType, argsCounter) + ") "
	} else {
		sq += column + operator + MakeParam(DBType, argsCounter) + " "
	}
	args = append(args, val)
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

func buildSQLTXTSearch(DBType byte, sq string, argsCounter int, operator string, val string, caseInsensitive bool, columns []string) (counter int, resquery string, args []interface{}) {
	suffix := " "
	suffixOr := " OR "
	suffixEnd := ") "
	if strings.Contains(val, `\`) || strings.Contains(val, "%") || strings.Contains(val, "_") {
		val = strings.Replace(val, `\`, `\\`, -1)
		val = strings.Replace(val, "%", `\%`, -1)
		val = strings.Replace(val, "_", `\_`, -1)
		suffix = ` ESCAPE '\' `
		suffixOr = ` ESCAPE '\' OR `
		suffixEnd = ` ESCAPE '\') `
	}
	val = "%" + val + "%"
	argsCounter++
	args = append(args, val)
	var valUpper string // will hold UPPERCASE string
	var valLower string // will hold lowercase string
	var valFirst string // will hold Firstletterupper string
	if caseInsensitive {
		valUpper = strings.ToUpper(val)
		valLower = strings.ToLower(val)
		valRune := []rune(val)
		ri := firstLetterIndex(valRune)
		valFirst = string(valRune[0:ri]) + strings.ToUpper(string(valRune[ri])) + strings.ToLower(string(valRune[ri+1:]))
		argsCounter += 3
		args = append(args, valUpper, valLower, valFirst)
	}
	for i := 0; i < len(columns); i++ {
		if strings.Contains(sq, "WHERE") {
			if i == 0 {
				sq += "AND ("
			} else {
				sq += "OR "
			}
		} else {
			sq += "WHERE ("
		}
		if caseInsensitive {
			sq += "(" + columns[i] + operator + MakeParam(DBType, argsCounter-3) + suffixOr +
				columns[i] + operator + MakeParam(DBType, argsCounter-2) + suffixOr +
				columns[i] + operator + MakeParam(DBType, argsCounter-1) + suffixOr +
				columns[i] + operator + MakeParam(DBType, argsCounter) + suffixEnd
		} else {
			sq += columns[i] + operator + MakeParam(DBType, argsCounter) + suffix
		}
	}
	sq += ") "
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

func buildUncountedSQLTXTSearch(DBType byte, sq string, argsCounter int, operator string, val string, caseInsensitive bool, columns []string) (counter int, resquery string, args []interface{}) {
	suffix := " "
	if strings.Contains(val, `\`) || strings.Contains(val, "%") || strings.Contains(val, "_") {
		val = strings.Replace(val, `\`, `\\`, -1)
		val = strings.Replace(val, "%", `\%`, -1)
		val = strings.Replace(val, "_", `\_`, -1)
		suffix = ` ESCAPE '\' `
	}
	val = "%" + val + "%"
	for i := 0; i < len(columns); i++ {
		argsCounter++
		if strings.Contains(sq, "WHERE") {
			if i == 0 {
				sq += "AND ("
			} else {
				sq += "OR "
			}
		} else {
			sq += "WHERE ("
		}
		sq += columns[i] + operator + MakeParam(DBType, argsCounter) + suffix
		if DBType == ORACLE && caseInsensitive {
			sq += "COLLATE binary_ai "
		}
		args = append(args, val)
	}
	sq += ") "
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

func buildSQLINJSONList(DBType byte, sq string, argsCounter int, column string, valueList []int) (counter int, resquery string, args []interface{}) {
	if strings.Contains(sq, "WHERE") {
		sq += "AND ("
	} else {
		sq += "WHERE ("
	}
	for i := 0; i < len(valueList); i++ {
		if i == 0 {
			sq += column + " LIKE " + MakeParam(DBType, argsCounter+1)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+2)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+3)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+4)
			argsCounter += 4
		} else {
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+1)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+2)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+3)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+4)
			argsCounter += 4
		}
		args = append(args, "%,"+strconv.Itoa(valueList[i])+",%")
		args = append(args, "%["+strconv.Itoa(valueList[i])+",%")
		args = append(args, "%,"+strconv.Itoa(valueList[i])+"]%")
		args = append(args, "%["+strconv.Itoa(valueList[i])+"]%")
	}
	sq += ") "
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

func buildSQLINJSONListOR(DBType byte, sq string, argsCounter int, column string, valueList []int, FirstIter bool, LastIter bool) (counter int, resquery string, args []interface{}) {
	if FirstIter {
		if strings.Contains(sq, "WHERE") {
			sq += "AND ("
		} else {
			sq += "WHERE ("
		}
	}
	if !FirstIter {
		sq += "OR ("
	}
	for i := 0; i < len(valueList); i++ {
		if i == 0 {
			sq += column + " LIKE " + MakeParam(DBType, argsCounter+1)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+2)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+3)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+4)
			argsCounter += 4
		} else {
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+1)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+2)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+3)
			sq += " OR " + column + " LIKE " + MakeParam(DBType, argsCounter+4)
			argsCounter += 4
		}
		args = append(args, "%,"+strconv.Itoa(valueList[i])+",%")
		args = append(args, "%["+strconv.Itoa(valueList[i])+",%")
		args = append(args, "%,"+strconv.Itoa(valueList[i])+"]%")
		args = append(args, "%["+strconv.Itoa(valueList[i])+"]%")
	}
	sq += ") "
	if LastIter {
		sq += ") "
	}
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

func buildSQLstrBETWEEN(DBType byte, sq string, argsCounter int, column string, valueList []int64) (counter int, resquery string, args []interface{}) {
	if strings.Contains(sq, "WHERE") {
		sq += "AND "
	} else {
		sq += "WHERE "
	}
	sq += column + " BETWEEN "
	for i := 0; i < len(valueList); i++ {
		argsCounter++
		if i == 0 {
			sq += MakeParam(DBType, argsCounter) + " "
		} else {
			sq += "AND " + MakeParam(DBType, argsCounter) + " "
		}
		args = append(args, valueList[i])
	}
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}

func buildSQLintBETWEEN(DBType byte, sq string, argsCounter int, column string, valueList []int) (counter int, resquery string, args []interface{}) {
	if strings.Contains(sq, "WHERE") {
		sq += "AND "
	} else {
		sq += "WHERE "
	}
	sq += column + " BETWEEN "
	for i := 0; i < len(valueList); i++ {
		argsCounter++
		if i == 0 {
			sq += MakeParam(DBType, argsCounter) + " "
		} else {
			sq += "AND " + MakeParam(DBType, argsCounter) + " "
		}
		args = append(args, valueList[i])
	}
	counter = argsCounter
	resquery = sq
	return counter, resquery, args
}
