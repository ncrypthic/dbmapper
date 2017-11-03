sqlmapper
=========

`sqlmapper` provides DSL for mapping columns from database into field/variable programmatically.
`sqlmapper` was inspired by Scala's [Anorm](https://www.playframework.com/documentation/latest/ScalaAnorm)
package. With mapping done on a per-column basis, adding new columns into existing table will not break
existing code as long as it is not modified.

Usage
-----

1. Import the package

   `import "github.com/ncrypthic/sqlmapper/dialects/cassandra"`

   or

   `import "github.com/ncrypthic/sqlmapper/dialects/mysql"`

2. Prepare result variables `result := make([]SomeStruct, 0)`

3. Create `RowMapper` interface instance, and use `MappedColumns.Then` to
   collect the result
   ```go
   // Prepare row to hold scan result
   func rowMapper() sqlmapper.RowMapper {
           return func() *sqlmapper.MappedColumns {
                   row := SomeStruct{}
                   sqlmapper.Columns(
                           sqlmapper.Column("column_name").As(&row.SomeField),
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
   mysql.Parse(sql.Query("SELECT * FROM some_table")).Map(rowMapper)
   ```

MySQL Example
----------------

```go
// ...
import (
    "database/sql"
    "github.com/ncrypthic/sqlmapper"
    "github.com/ncrypthic/sqlmapper/dialect/mysql"
)

type Dummy struct {
    ID       string
    Name     string
    Active   bool
    OptField sql.NullString
}

func main() {
        // db := sql.Open()
        result := make([]Dummy, 0)

        // Single row mapper for table 'dummy'. This way it can be reused anywhere
        dummySqlMapper := func(result *[]Dummy) *sqlmapper.MappedColumns {
                row := Dummy{}
                return sqlmapper.Columns(
                        sqlmapper.Column("id").As(&row.ID),
                        sqlmapper.Column("id").As(&row.Name),
                        sqlmapper.Column("id").As(&row.Active),
                        sqlmapper.Column("id").As(&row.OptField),
                ).Then(func() error {
                        *result = append(*result, row)
                        return nil
                })
        }

        // Multiple rows mapper for table 'dummy'
        dummyRowsMapper = func(result *[]Dummy) sqlmapper.RowMapper {
                return func() *sqlmapper.MappedColumns {
                        return dummySqlMapper(result)
                }
        }

        sqlmapper.Parse(db.Query("SELECT id, name, active, opt_field FROM example")).Map(dummyRowsMapper(result))

        for _, r := range result {
            // Do something with the result
        }
}
```

LICENSE
-------

MIT Licensed
