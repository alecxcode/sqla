// Package sqla provides some specific functions to extend standard Go database/sql package.
// These functions can be used in any SQL-driven applications, although the pachage initially was created for the EDM project
// [see https://github.com/alecxcode/edm for the most complete example of sqla usage].
//
// This package is intended to provide more convenient methods for accessing SQL databases:
// creating, updating, deleting and selecting objects.
// Standard Go database/sql functions are not changed.
// All new functions works with them, and usual database/sql should be used when necessary.
//
// The package supports the following RDBMS: SQLite, Microsoft SQL Server, MySQL(MariaDB), Oracle, PostgreSQL.
//
// The key functions of this package are related to the following:
// working with different RDBMS seamlessly;
// constructing select statement with multiple filters programmatically and arguments list protected from SQL injection;
// easier (than with bare database/sql) inserting, updating, deleting objects.
//
package sqla
