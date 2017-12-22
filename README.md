dbmapper
========

`dbmapper` provides DSL for mapping database result set into specific field/variable `programmatically` and it also provide DSL to help developer create a query with named parameter. `dbmapper` was inspired by Scala's [Anorm](https://www.playframework.com/documentation/latest/ScalaAnorm) package. With mapping done on a per-column basis, adding new columns into existing table will not break existing code as long as it is not modified.

Query Usage
===========

1. Import the package

   `import "github.com/ncrypthic/dbmapper"`

2. Import database dialect. Currently only `mysql` and `cassandra` are supported
   ```go
   // import "github.com/ncrypthic/dbmapper/dialects/cassandra"
   // or
   import "github.com/ncrypthic/dbmapper/dialects/mysql"
   ```

3. Create database query with named parameter
   ```go
   query := "SELECT col_a, col_b FROM a_table WHERE col_a = :a_parameter)
   ```

4. You can convert the parameterized query to native driver query using `dbmapper.Prepare` api
   ```go
   queryString := "SELECT col_a, col_b FROM a_table WHERE col_a = :a_parameter)
   query := dbmapper.Prepare(queryString).With(
           mysql.Param("a_parameter", "some_value"),
   )
   resultSet, err := driver.Query(query.SQL(), query.Params()...)
   ```
   or
   ```go
   queryString := "SELECT col_a, col_b FROM a_table WHERE col_a = :a_parameter)
   query := dbmapper.Prepare(queryString).With(
           cassandra.Param("a_parameter", "some_value"),
   )
   resultSet, err := driver.Query(query.SQL(), query.Params()...)
   ```

Result Mapping Usage
====================

1. Import the package

   `import "github.com/ncrypthic/dbmapper"`

2. Import database dialect. Currently only `mysql` and `cassandra` are supported
   ```go
   import "github.com/ncrypthic/dbmapper/dialects/cassandra"
   ```
   or
   ```go
   import "github.com/ncrypthic/dbmapper/dialects/mysql"
   ```

2. Prepare variables to hold the mapped row results `result := make([]SomeStruct, 0)`

3. Create `RowMapper` interface instance, and use `MappedColumns.Then` to
   collect the result
   ```go

   func rowMapper(result []SomeStruct) dbmapper.RowMapper {
           return func() *dbmapper.MappedColumns {
                   // Either create new object to hold row result or use existing
                   row := SomeStruct{}
                   dbmapper.Columns(
                           dbmapper.Column("column_name").As(&row.SomeField),
                   ).Then(func() error {
                           // Append the row to a result slice
                           result = append(result, row)
                           // When error returned, it will stop the mapping process
                           return nil
                   })
           }
   }
   ```


5. Pass query result from native database driver to `<dialect_package>.Parse` method then pass instance of `RowMapper` to map the result
   ```go
   mysql.Parse(sql.Query("SELECT col_1, col_2 FROM some_table")).Map(rowMapper(result))
   ```

   Or

   ```go
   cassandra.Parse(gocql.Query("SELECT col_1, col_2 FROM some_table")).Map(rowMapper(result))
   ```

Example
=======

```go
// ...
import (
    "database/sql"
    "github.com/ncrypthic/dbmapper"
    "github.com/ncrypthic/dbmapper/dialects/mysql"
)

type User struct {
    ID          string
    Name        string
    Active      bool
    PhoneNumber sql.NullString
}

// userMapper map a single row from table 'user'. This way it can be reused anywhere
func userMapper(result []User) dbmapper.RowMapper {
        return func() *dbmapper.MappedColumns {
                row := User{}
                return dbmapper.Columns(
                        dbmapper.Column("id").As(&row.ID),
                        dbmapper.Column("name").As(&row.Name),
                        dbmapper.Column("active").As(&row.Active),
                        dbmapper.Column("phone_number").As(&row.OptField),
                ).Then(func() error {
                        result = append(result, row)
                        return nil
                })
        }
}

// usersMapper map rows of query result from table 'user' using `userMapper` function to map every row
func usersMapper(result []User) dbmapper.RowMapper {
        return userMapper(result)
}

func main() {
        // db := sql.Open()
        sql := "SELECT id, name, active, phone_number FROM user"
        users := make([]User, 0)

        query := dbmapper.Prepare("SELECT id, name, active, phone_number FROM user")

        dbmapper.Parse(db.Query(query.SQL(), query.Params()...)).Map(usersMapper(result))

        for _, r := range result {
            // Do something with the result
        }
}
```

LICENSE
-------

MIT
