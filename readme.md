# sqla Go (Golang) package

Package **sqla** provides some specific functions to extend standard Go database/sql package. These functions can be used in any SQL-driven applications, although the package initially was created for the EDM project [see https://github.com/alecxcode/edm for the most complete example of sqla usage]. Basic example code (how to use) can be found below.

Complete documentation available at: https://pkg.go.dev/github.com/alecxcode/sqla

This package is intended to provide more convenient methods for accessing SQL databases: creating, updating, deleting and selecting objects.
Standard Go database/sql functions are not changed. All new functions works with them, and usual database/sql should be used when necessary.

## Features and functions:

### The package supports the following RDBMS:
* SQLite
* Microsoft SQL Server
* MySQL(MariaDB)
* Oracle
* PostgreSQL

### The key functions of this package are related to the following:
* working with different RDBMS seamlessly;
* constructing select statement with multiple filters programmatically and arguments list protected from SQL injection;
* easier (than with bare database/sql) inserting, updating, deleting objects.

## Installation and use:

### How to use in your project:

Add `"github.com/alecxcode/sqla"` to you import section in a *.go file.

Run in the package folder:  
```shell
go mod init nameofyourproject  
go mod tidy  
```

### Basic use code example:

```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/alecxcode/sqla"
)

// Book is a type to represent books
type Book struct {
	ID            int
	BookTitle     string
	Author        string
	YearPublished int
}

func (b *Book) create(db *sql.DB, DBType byte) (lastid int, rowsaff int) {
	var args sqla.AnyTslice
	args = args.AppendNonEmptyString("BookTitle", b.BookTitle)
	args = args.AppendNonEmptyString("Author", b.Author)
	args = args.AppendInt("YearPublished", b.YearPublished)
	lastid, rowsaff = sqla.InsertObject(db, DBType, "books", args)
	return lastid, rowsaff
}
func (b *Book) update(db *sql.DB, DBType byte) (rowsaff int) {
	var args sqla.AnyTslice
	args = args.AppendStringOrNil("BookTitle", b.BookTitle)
	args = args.AppendStringOrNil("Author", b.Author)
	args = args.AppendInt("YearPublished", b.YearPublished)
	rowsaff = sqla.UpdateObject(db, DBType, "books", args, b.ID)
	return rowsaff
}
func (b *Book) load(db *sql.DB, DBType byte) error {
	row := db.QueryRow(
		"SELECT ID, BookTitle, Author, YearPublished FROM books WHERE ID = "+
			sqla.MakeParam(DBType, 1),
		b.ID)
	var BookTitle, Author sql.NullString
	var YearPublished sql.NullInt64
	err := row.Scan(&b.ID, &BookTitle, &Author, &YearPublished)
	if err != nil {
		return err
	}
	b.BookTitle = BookTitle.String
	b.Author = Author.String
	b.YearPublished = int(YearPublished.Int64)
	return nil
}

// App is a type to represent computer software
type App struct {
	ID           int
	AppName      string
	Author       string
	YearReleased int
}

func (a *App) create(db *sql.DB, DBType byte) (lastid int, rowsaff int) {
	var args sqla.AnyTslice
	args = args.AppendNonEmptyString("AppName", a.AppName)
	args = args.AppendNonEmptyString("Author", a.Author)
	args = args.AppendInt("YearReleased", a.YearReleased)
	lastid, rowsaff = sqla.InsertObject(db, DBType, "apps", args)
	return lastid, rowsaff
}
func (a *App) update(db *sql.DB, DBType byte) (rowsaff int) {
	var args sqla.AnyTslice
	args = args.AppendStringOrNil("AppName", a.AppName)
	args = args.AppendStringOrNil("Author", a.Author)
	args = args.AppendInt("YearReleased", a.YearReleased)
	rowsaff = sqla.UpdateObject(db, DBType, "apps", args, a.ID)
	return rowsaff
}
func (a *App) load(db *sql.DB, DBType byte) error {
	row := db.QueryRow(
		"SELECT ID, AppName, Author, YearReleased FROM apps WHERE ID = "+
			sqla.MakeParam(DBType, 1),
		a.ID)
	var AppName, Author sql.NullString
	var YearReleased sql.NullInt64
	err := row.Scan(&a.ID, &AppName, &Author, &YearReleased)
	if err != nil {
		return err
	}
	a.AppName = AppName.String
	a.Author = Author.String
	a.YearReleased = int(YearReleased.Int64)
	return nil
}

func main() {

	// Initializing database
	const DBType = sqla.SQLITE
	var db *sql.DB
	db = sqla.OpenSQLConnection(DBType, "file::memory:?cache=shared&_foreign_keys=true")
	db.Exec("CREATE TABLE books (ID INTEGER PRIMARY KEY, BookTitle TEXT, Author TEXT, YearPublished INTEGER);")
	db.Exec("CREATE TABLE apps (ID INTEGER PRIMARY KEY, AppName TEXT, Author TEXT, YearReleased INTEGER);")

	// Creating objects and inserting into database
	var someBook = Book{BookTitle: "Alice's Adventures in Wonderland", Author: "Lewis Carroll", YearPublished: 1865}
	var someApp = App{AppName: "Linux", Author: "Linus Torvalds", YearReleased: 1991}
	bookID, res := someBook.create(db, DBType)
	if res > 0 {
		fmt.Println("Inserted book into DB")
		someBook.ID = bookID
	}
	appID, res := someApp.create(db, DBType)
	if res > 0 {
		fmt.Println("Inserted app into DB")
		someApp.ID = appID
	}

	// Updating object in the database
	someBook.BookTitle = "Some Updated Book Title"
	someBook.Author = ""
	someBook.YearPublished = 1900
	res = someBook.update(db, DBType)
	if res > 0 {
		fmt.Println("Updated book in the DB")
	}

	// Loading objects from database
	bookFromDB := Book{ID: bookID}
	appFromDB := App{ID: appID}
	bookFromDB.load(db, DBType)
	appFromDB.load(db, DBType)
	fmt.Printf("Book loaded from DB: %#v\n", bookFromDB)
	fmt.Printf("App loaded from DB: %#v\n", appFromDB)

	// Deleting objects from database
	res = sqla.DeleteObjects(db, DBType, "books", "ID", []int{bookFromDB.ID})
	if res > 0 {
		fmt.Println("Deleted book from DB")
	}
	res = sqla.DeleteObjects(db, DBType, "apps", "ID", []int{appFromDB.ID})
	if res > 0 {
		fmt.Println("Deleted app from DB")
	}

	db.Close()

}
```