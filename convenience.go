package sqlconvenient

import (
	"database/sql"
	"log"
	"reflect"
)

/*
Prepares the query string by inserting given parameters
and then executes the query.
For example used for SQL insert operations.
*/
func SqlExec(db *sql.DB, query string, params ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		log.Fatal(err)
		return err
	}

	tx.Commit()
	return nil
}

/*
Passes params into query string, executes the query and returns a slice
of type rettype disguised as interface.
This is useful when the real type of the slice isn't actually needed. For
example passing the result into a json encoder works fine even with this
interface slice.
If need be, one can convert the results in 3 lines into a slice with correct
type:
	realslice := make([]REALTYPE, len(results))
	for i, val := range results {
		realslice[i] = val.(REALTYPE)
	}

This is used for SQL select operations.
rettype can not only be a struct, but also a simple type such as string or int.

Example usage:
	var p Product
	result := sqlQueryStruct(&p, "...")

B E W A R E
This uses reflect, thus is inefficient and purely for convenience.
*/
func SqlQuery(db *sql.DB, rettype interface{}, query string, params ...interface{}) []interface{} {
	var ret []interface{}

	stmt, _ := db.Prepare(query)
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		log.Fatal(err)
	}

	isStruct := reflect.ValueOf(rettype).Elem().Kind() == reflect.Struct

	for rows.Next() {
		var ptr []interface{}

		// Create a new object of type rettype
		val := reflect.ValueOf(rettype).Elem()

		// If it's a struct, add pointers to all its fields into ptr slice,
		// else add pointer to the object itself
		if isStruct {
			for i := 0; i < val.NumField(); i++ {
				ptr = append(ptr, val.Field(i).Addr().Interface())
			}
		} else {
			ptr = append(ptr, val.Addr().Interface())
		}

		// Pass pointers into scan which fills them
		err = rows.Scan(ptr...)
		if err != nil {
			log.Fatal(err)
		}

		ret = append(ret, val.Interface())
	}

	return ret
}
