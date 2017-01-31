dbsql2go
========

dbsql2go generates Go structs for all tables in a given database. If the column is nullable, the appropriate `sql.Null%` type will be used, if applicable.  If the database driver package implements additional `sql.Null` types, those will also be used when appropriate. This does not apply to types that are either binary or resolve to []byte, The Go type used will be the type that most closely matches the db column's type. Support for some db specfic types may be missing and may be implemented in the future.

This is meant to take the initial busy work out of using a relational database without resorting to an ORM. The generated functions and methods are limited.

In addition to generating structs, DML methods and functions will be generated as appropriate.

Currently, any table with a primary key will have pk based `SELECT`, `UPDATE`. and `DELETE` methods for single row operations. Range `SELECT` funcs will also be generated for multiple row operation.

All tables will have an `INSERT` method defined.

Views only have structs defined for them.

It is assumed that the login user used has the necessary permissions to query the RDBMSs database catalogs.

## Usage

To generate Go code for the `dbname` MySQL database:

    $ dbsql2go -rdbms mysql -db dbname -user dbuser -password notapassword

## Flags
Not all flags are required. For `user` and `password` either the long flag or the short flag is required.

Flag | Type | Default | Required | Description  
:--|:--|:--:|:--:|:--  
rdbms|string||true|The target RDBMS  
db|string||true|Database name  
user|string||true|Login user  
u|string||true|Login ser (short)  
password|string||true|User's password  
p|string||true|User's password (short)  
server|string||RDBMs dependent|Server location
package|string||false|Name of the package of which the generated code is a part; if empty thje WD name will be used  
dbpackage|bool|false|false|Use the database name as the package name; this overrides the package string  
filepertable|bool|false|false|use a file per table  

## Supported databases

### MySQL
The MySQL driver is [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).

Support for geo types is not implemented; any columns using one of these types currently has the MySQL type as it's Go type, which is clearly wrong. These will need to be replaced by the user until support for those types has been added.

The user must have `SELECT` permissions on the `information_schema`.
