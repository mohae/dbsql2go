dbsql2go
========

dbsql2go generates Go structs for all tables in a given database. If the column is nullable, the appropriate `sql.Null%` type will be used, if applicable. This does not apply to types that are either binary or resolve to []byte, The Go type used will be the type that most closely matches the db column's type. Support for some db specfic types may be missing and may be implemented in the future.

It is assumed that the login user used has the necessary permissions to query the RDBMSs database catalogs.


## Supported databases

### MySQL
The MySQL driver is [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).

Support for geo types is not implemented; any columns using one of these types currently has the MySQL type as it's Go type, which is clearly wrong. These will need to be replaced by the user until support for those types has been added.

The user must have `SELECT` permissions on the `information_schema`.


