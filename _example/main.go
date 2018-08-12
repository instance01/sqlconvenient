package main

import (
	"database/sql"
	"fmt"
	sqlconvenient "github.com/instance01/sqlconvenient"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println("Error: Could not initialize database. Does test.db exist?")
	}
	return db, err
}

type Entry struct {
	Foo int
	Bar string
}

func getAllEntries() []interface{} {
	var e Entry
	var query = "select * from testtable"
	return sqlconvenient.SqlQuery(db, &e, query)
}

func getGroupNames(id int) []interface{} {
	var s string
	var query = "select bar from testtable where foo = ? group by bar"
	return sqlconvenient.SqlQuery(db, &s, query, id)
}

func insertEntry(id int, name string) error {
	var query = "insert into testtable (foo, bar) values(?, ?)"
	return sqlconvenient.SqlExec(db, query, id, name)
}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		return
	}

	insertEntry(3, "test222")

	entries := getAllEntries()
	fmt.Printf("%#v\n", entries)

	names := getGroupNames(3)
	fmt.Printf("%s\n%s\n", names...)
}
