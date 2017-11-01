sqlmapper
=========

sqlmapper provides DSL for mapping columns into field/variable programmatically.
With controlled mapping, adding new columns into existing table in database
will not break existing code as long as existing column is not modified.

Usage
-----

1. Import the package

`import "github.com/ncrypthic/sqlmapper"

2. Prepare result variables `result := make([]SomeStruct, 0)`

3. Create `RowMapper` interface instance, and use `MappedColumns.Then` to
   collect the result
   ```go
           row := SomeStruct{}
           sqlmapper.Columns(
                   sqlmapper.Column("column_name").As(&row.SomeField),
           ).Then(func() error {
                   // collect the result
                   result = append(result, row)
                   return nil
           })
   ```

4. Pass `sql.Query` return value into `sqlmapper.Parse` method then
   pass Row mapper to map the result
   ```
   sqlmapper.Parse(sql.Query("SELECT * FROM some_table")).Map(rowMapperInstance)
   ```

Complete Example
----------------

```go
// ...
import (
    "database/sql"
    "github.com/ncrypthic/sqlmapper"
)

type Dummy struct {
    ID       string
    Name     string
    Active   bool
    OptField sql.NullString
}

// db := sql.Open()
result := make([]Dummy, 0)

// Single row mapper for table 'dummy'. This way it can be reused anywhere
dummySqlMapper := func(result *[]Dummy) MappedColumns {
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
dummyRowsMapper = func(result *[]Dummy) RowMapper {
        return func() MappedColumns {
                return dummySqlMapper(result)
        }
}

sqlmapper.Parse(db.Query("SELECT id, name, active, opt_field FROM example")).Map(func() MappedColumns {
        return dummySqlMapper(row)
})
// OR
sqlmapper.Parse(db.Query("SELECT id, name, active, opt_field FROM example")).Map(dummyRowsMapper(result))

for _, r := range result {
    // Do something with the result
}
```

LICENSE
-------

MIT Licensed
