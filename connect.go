package sqla

import (
	"database/sql"
	"log"
	"net/url"
	"strings"
	"time"

	//names of drivers omitted
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora/v2"
	//_ "github.com/lib/pq"
)

// SQLITE, MSSQL, MYSQL, ORACLE POSTGRESQL - are database types supported.
const (
	SQLITE = iota
	MSSQL
	MYSQL
	ORACLE
	POSTGRESQL
)

// DEBUG may be set to true to print SQL queries
const DEBUG = false

// ReturnDBType gives digital representation of RDBMS type based on string.
// Accepted db types are: "sqlite", "mssql" (or "sqlserver"), "mysql" (or "mariadb"), "oracle", "postgresql" (or "postgres").
func ReturnDBType(dbtype string) byte {
	switch dbtype {
	case "sqlite":
		return SQLITE
	case "mssql":
		return MSSQL
	case "sqlserver":
		return MSSQL
	case "mysql":
		return MYSQL
	case "mariadb":
		return MYSQL
	case "oracle":
		return ORACLE
	case "postgres":
		return POSTGRESQL
	case "postgresql":
		return POSTGRESQL
	}
	return 0
}

// BuildDSN creates DSN for database connection. Then, DSN is used in CreateDB and OpenSQLConnection.
// Accepted db types are: "sqlite", "mssql" (or "sqlserver"), "mysql" (or "mariadb"), "oracle", "postgresql" (or "postgres").
// Other arguments are self-descripting. For Oracle DBName is a service name.
func BuildDSN(DBType string, DBName string, DBHost string, DBPort string, DBUser string, DBPassword string) (DSN string) {
	DBTypeC := ReturnDBType(DBType)
	if DBTypeC == SQLITE {
		DSN = DBName
	} else if DBTypeC == MSSQL {
		//DSN = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;",
		//	DBHost, DBUser, DBPassword, DBPort, DBName)
		DSN = "sqlserver://" + url.QueryEscape(DBUser) + ":" + url.QueryEscape(DBPassword) + "@" + DBHost + ":" + DBPort + "?database=" + DBName + "&encrypt=disable"
	} else if DBTypeC == MYSQL {
		// This MySQL driver developers say no need for escaping
		DSN = DBUser + ":" + DBPassword + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName
	} else if DBTypeC == ORACLE {
		DSN = "oracle://" + url.QueryEscape(DBUser) + ":" + url.QueryEscape(DBPassword) + "@" + DBHost + ":" + DBPort + "/" + DBName
	} else if DBTypeC == POSTGRESQL {
		//DSN = "host=" + DBHost + " dbname=" + DBName + " user=" + DBUser + " password=" + DBPassword + " port=" + DBPort + " sslmode=disable"
		DSN = "postgres://" + url.QueryEscape(DBUser) + ":" + url.QueryEscape(DBPassword) + "@" + DBHost + ":" + DBPort + "/" + DBName + "?sslmode=disable"
	}
	//log.Println(DSN)
	return DSN
}

// CreateDB creates database for SQLITE and schema for all databases, based on provided sql script (sqlStmt argument).
// CreateDB automatically opens database connection and then closes the connection after creation is complete.
// For DBType see constants, for DSN see BuildDSN.
func CreateDB(DBType byte, DSN string, sqlStmt string) {

	var err error
	db := OpenSQLConnection(DBType, DSN)
	defer db.Close()

	sqlStmtArr := strings.Split(sqlStmt, ";")
	for i := 0; i < len(sqlStmtArr)-1; i++ {
		_, err = db.Exec(strings.Trim(sqlStmtArr[i], "\r\n\t ;"))
		if err != nil {
			log.Println(strings.Trim(sqlStmtArr[i], "\r\n\t ;"))
			log.Printf("%q: %s%d\n", err, "while creating tables at:", i)
		}
	}

}

// OpenSQLConnection onpens connection to a database and renurns standard Go *sql.DB type.
// For DBType see constants, for DSN see BuildDSN.
func OpenSQLConnection(DBType byte, DSN string) (db *sql.DB) {
	var err error
	var sqldriver string
	switch DBType {
	case SQLITE:
		sqldriver = "sqlite3"
	case MSSQL:
		sqldriver = "sqlserver"
	case MYSQL:
		sqldriver = "mysql"
	case ORACLE:
		sqldriver = "oracle"
	case POSTGRESQL:
		sqldriver = "postgres"
	}
	db, err = sql.Open(sqldriver, DSN)
	if err != nil {
		log.Fatal("Opening SQL connection:", err)
	}
	if DBType == MYSQL {
		db.SetConnMaxLifetime(time.Minute * 3) // MySQL closes connection on idle, we close first to avoid errors
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Opening SQL connection:", err)
	}
	return db
}
