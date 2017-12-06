dbmapper
========

`dbmapper` provides DSL for mapping columns from database into field/variable programmatically.
`dbmapper` was inspired by Scala's [Anorm](https://www.playframework.com/documentation/latest/ScalaAnorm)
package. With mapping done on a per-column basis, adding new columns into existing table will not break
existing code as long as it is not modified.

Usage
-----

1. Import the package

   `import "github.com/ncrypthic/dbmapper/dialects/cassandra"`

   or

   `import "github.com/ncrypthic/dbmapper/dialects/mysql"`

2. Prepare result variables `result := make([]SomeStruct, 0)`

3. Create `RowMapper` interface instance, and use `MappedColumns.Then` to
   collect the result
   ```go
   // Prepare row to hold scan result
   func rowMapper() dbmapper.RowMapper {
           return func() *dbmapper.MappedColumns {
                   row := SomeStruct{}
                   dbmapper.Columns(
                           dbmapper.Column("column_name").As(&row.SomeField),
                   ).Then(func() error {
                           // append row to a result slice
                           // result = append(result, row)
                           return nil
                   })
           }
   }
   ```

4. Pass `sql.Query` / `*gocql.Query` return value into `<dialect_package>.Parse` method then
   pass Row mapper to map the result
   ```
   mysql.Parse(sql.Query("SELECT col_1, col_2 FROM some_table")).Map(rowMapper)
   ```

   Or for cassandra

   ```
   cassandra.Parse(gocql.Query("SELECT col_1, col_2 FROM some_table")).Map(rowMapper)
   ```

Usage example
-------------

```go
// ...
import (
    "database/sql"
    "github.com/ncrypthic/dbmapper"
    "github.com/ncrypthic/dbmapper/dialect/mysql"
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
        result := make([]User, 0)

        dbmapper.Parse(db.Query(sql)).Map(usersMapper(result))

        for _, r := range result {
            // Do something with the result
        }
}
```

LICENSE
-------

MIT
