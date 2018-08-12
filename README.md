
This package exports two functions

```
SqlExec(db *sql.DB, query string, params ...interface{})

SqlQuery(db *sql.DB, rettype interface{}, query string, params ...interface{})
```

and allows querying/executing sql operations in just 2-3 lines, resulting in less boilerplate. Example:

```
var e Entry
var query = "select * from testtable"
return sqlconvenient.SqlQuery(db, &e, query)
```

See a more in depth example in \_example/ and documentation [over here](https://godoc.org/github.com/instance01/sqlconvenient).
