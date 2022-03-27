package sqla

import (
	"log"
	"strings"
)

// Seek type allows to implement so-called seek method of pagination.
// If UseSeek is true, then ConstructSELECTquery will use seek method (instead of offsetting) which just selects values greater or less than Seek.Value.
//
// If ValueInclude is true, then seek method uses greater or equal (or less or equal) operator. This behavior may be useful when reloading the same page.
type Seek struct {
	UseSeek      bool
	Value        int
	ValueInclude bool
}

// ConstructSELECTquery is the most fundamental for this package. It constructs and returns SQL statement for select query, SQL statement for COUNT(), and arguments slice to use in Go sql functions.
// The function takes the following arguments: DBType (see constants);
// tableName to paste after FROM; columnsToSelect as a comma-separated columns;
// columnsToCount - to put as an argument for a COUNT(), e.g. "*" will be COUNT(*);
// joins as usual joins part of an SQL statement;
// Filter - is the main thing to counstruct query based on different filters. See Filter type and its methods;
// orderBy is a column name or column's names to order result; limit, offset - are usual values for sql statement;
// distinct as bool defines whether you need to add DISTINCT keyword in your statement.
//
// Seek is used to avoid offsetting when dealing with big tables and to implement so-called seek method of pagination. See Seek type.
// Seek method of pagination requires additional coding in your app, and algorithms are not so simple as with offset.
func ConstructSELECTquery(
	DBType byte,
	tableName string,
	columnsToSelect string,
	columnsToCount string,
	joins string,
	F Filter,
	orderBy string,
	orderHow int,
	limit int,
	offset int,
	distinct bool,
	seek Seek) (sq string, sqcount string, args []interface{}, argscount []interface{}) {

	var argstoAppend []interface{}
	var argsCounter int

	for _, FC := range F.ClassFilter {
		if FC.InJSON {
			argsCounter, sq, argstoAppend = buildSQLINJSONList(DBType, sq, argsCounter, FC.Column, FC.List)
			args = append(args, argstoAppend...)
		} else {
			argsCounter, sq, argstoAppend = BuildSQLIN(DBType, sq, argsCounter, FC.Column, FC.List)
			args = append(args, argstoAppend...)
		}
	}

	var CurrentFCName string
	if len(F.ClassFilterOR) > 0 {
		CurrentFCName = F.ClassFilterOR[0].Name
	}
	for i, FC := range F.ClassFilterOR {
		FirstIter := false
		LastIter := false
		if i == 0 {
			FirstIter = true
		}
		if i == len(F.ClassFilterOR)-1 {
			LastIter = true
		}
		if i+1 < len(F.ClassFilterOR) {
			if FC.Name != F.ClassFilterOR[i+1].Name {
				LastIter = true
			}
		}
		if CurrentFCName != FC.Name {
			FirstIter = true
			CurrentFCName = FC.Name
		}
		if FC.InJSON {
			argsCounter, sq, argstoAppend = buildSQLINJSONListOR(DBType, sq, argsCounter, FC.Column, FC.List, FirstIter, LastIter)
			args = append(args, argstoAppend...)
		} else {
			argsCounter, sq, argstoAppend = BuildSQLINOR(DBType, sq, argsCounter, FC.Column, FC.List, FirstIter, LastIter)
			args = append(args, argstoAppend...)
		}
	}

	for _, DF := range F.DateFilter {
		if len(DF.Dates) > 0 {
			operator := getRelationFromString(DF.Relation)
			if len(DF.Dates) == 1 {
				argsCounter, sq, argstoAppend = buildSQLCOMPARE(DBType, sq, argsCounter, DF.Column, operator, DF.Dates[0])
				args = append(args, argstoAppend...)
			} else {
				argsCounter, sq, argstoAppend = buildSQLstrBETWEEN(DBType, sq, argsCounter, DF.Column, DF.Dates)
				args = append(args, argstoAppend...)
			}
		}
	}

	for _, SF := range F.SumFilter {
		if len(SF.Sums) > 0 {
			operator := getRelationFromString(SF.Relation)
			if len(SF.Sums) == 1 {
				argsCounter, sq, argstoAppend = buildSQLCOMPARE(DBType, sq, argsCounter, SF.Column, operator, SF.Sums[0])
				args = append(args, argstoAppend...)
			} else {
				argsCounter, sq, argstoAppend = buildSQLintBETWEEN(DBType, sq, argsCounter, SF.Column, SF.Sums)
				args = append(args, argstoAppend...)
			}
			if SF.CurrencyCode != 0 {
				argsCounter, sq, argstoAppend = buildSQLCOMPARE(DBType, sq, argsCounter, SF.CurrencyColumn, " = ", SF.CurrencyCode)
				args = append(args, argstoAppend...)
			}
		}
	}

	if F.TextFilter != "" {
		operator := " LIKE "
		caseins := false
		if DBType == POSTGRESQL {
			operator = " ILIKE "
		}
		if (DBType == SQLITE && !isStringASCII(F.TextFilter)) || DBType == ORACLE {
			caseins = true
		}
		if DBType == MYSQL || DBType == ORACLE {
			// MySQL, Oracle, and others with ? or unaccessible by number placeholder:
			argsCounter, sq, argstoAppend = buildUncountedSQLTXTSearch(DBType, sq, argsCounter, operator, F.TextFilter, caseins, F.TextFilterColumns)
			args = append(args, argstoAppend...)
		} else {
			argsCounter, sq, argstoAppend = buildSQLTXTSearch(DBType, sq, argsCounter, operator, F.TextFilter, caseins, F.TextFilterColumns)
			args = append(args, argstoAppend...)
		}
	}

	if distinct {
		sqcount = "SELECT COUNT(DISTINCT " + columnsToCount + ") FROM " + tableName + " " + joins + " " + sq
		sq = "SELECT DISTINCT " + columnsToSelect + " FROM " + tableName + " " + joins + " " + sq
	} else {
		sqcount = "SELECT COUNT(" + columnsToCount + ") FROM " + tableName + " " + joins + " " + sq
		sq = "SELECT " + columnsToSelect + " FROM " + tableName + " " + joins + " " + sq
	}

	argscount = make([]interface{}, len(args))
	copy(argscount, args)

	if seek.UseSeek {
		operator := " > "
		if orderHow == 0 && !seek.ValueInclude {
			operator = " < "
		} else if orderHow == 0 && seek.ValueInclude {
			operator = " <= "
		} else if orderHow == 1 && seek.ValueInclude {
			operator = " >= "
		}
		argsCounter, sq, argstoAppend = buildSQLCOMPARE(DBType, sq, argsCounter, orderBy, operator, seek.Value)
		args = append(args, argstoAppend...)
	}

	var ordarr []string
	ordarr = strings.Split(orderBy, ",")

	if len(orderBy) > 0 {
		sq += "ORDER BY "
		for i, col := range ordarr {
			sq += col
			if DBType == SQLITE {
				//TODO?: add similar case-insensitive ORDER BY options for other RDBMS
				sq += " COLLATE NOCASE"
			}
			if orderHow == 0 {
				sq += " DESC "
			} else {
				sq += " ASC "
			}
			if (i + 1) < len(ordarr) {
				sq += ", "
			}
		}
	}

	if DBType == MSSQL || DBType == ORACLE {
		argsCounter++
		sq += " OFFSET " + MakeParam(DBType, argsCounter)
		argsCounter++
		sq += " ROWS FETCH NEXT " + MakeParam(DBType, argsCounter) + " ROWS ONLY"
		args = append(args, offset, limit)
	} else {
		argsCounter++
		sq += " LIMIT " + MakeParam(DBType, argsCounter)
		argsCounter++
		sq += " OFFSET " + MakeParam(DBType, argsCounter)
		args = append(args, limit, offset)
	}

	if DEBUG {
		log.Println(sq, args)
		log.Println(sqcount, argscount)
	}
	return sq, sqcount, args, argscount
}
