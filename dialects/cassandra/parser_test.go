package cassandra

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	. "github.com/ncrypthic/sqlmapper"
)

func mockRow(id uint, name string, active bool, opt_field *string) []interface{} {
	return []interface{}{id, name, active, opt_field}
}

type mockQuery struct {
}

func (i *mockQuery) Iter() CqlIterator {
	return &mockMapScaner{}
}

type mockMapScaner struct {
	pos int
}

func (i *mockMapScaner) Columns() []gocql.ColumnInfo {
	return []gocql.ColumnInfo{
		gocql.ColumnInfo{Name: "id"},
		gocql.ColumnInfo{Name: "name"},
		gocql.ColumnInfo{Name: "active"},
		gocql.ColumnInfo{Name: "opt_field"},
	}
}

func (i *mockMapScaner) Scan(result ...interface{}) bool {
	alice := "alice"
	pAlice := &alice
	bob := "bob"
	pBob := &bob
	charlie := "charlie"
	rows := [][]interface{}{
		mockRow(1, alice, true, pAlice),
		mockRow(2, bob, false, pBob),
		mockRow(2, charlie, false, nil),
	}
	if i.pos >= len(rows) {
		return false
	} else {
		for j, v := range rows[i.pos] {
			i.assignVal(result[j], v)
		}
		i.pos++
		return true
	}
}

func (i *mockMapScaner) assignVal(dest interface{}, source interface{}) {
	switch d := dest.(type) {
	case **string:
		switch s := source.(type) {
		case *string:
			*d = s
		}
	case *string:
		switch s := source.(type) {
		case string:
			*d = s
		}
	case *int:
		switch s := source.(type) {
		case int:
			*d = s
		}
	case *bool:
		switch s := source.(type) {
		case bool:
			*d = s
		}
	}
}

type ParseErr error

type User struct {
	ID        string
	Name      string
	Active    bool
	OptString *string
}

func userSqlMapper(result *[]User) *MappedColumns {
	user := User{ID: "1"}
	return Columns(
		Column("id").As(&user.ID),
		Column("name").As(&user.Name),
		Column("active").As(&user.Active),
		Column("opt_field").As(&user.OptString),
	).Then(func() error {
		*result = append(*result, user)
		return nil
	})
}

func usersSqlMapper(result *[]User) RowMapper {
	return func() *MappedColumns {
		return userSqlMapper(result)
	}
}

func Query(query string) CqlQuery {
	return &mockQuery{}
}

func TestRowParser(t *testing.T) {
	query := "SELECT id, name, active, opt_string FROM users"
	users := make([]User, 0)
	err := Parse(Query(query)).Map(func() *MappedColumns {
		return userSqlMapper(&users)
	})
	if err != nil {
		fmt.Printf("%v", ParseErr(err))
	} else {
		fmt.Println("result:")
		for _, r := range users {
			fmt.Printf("%+v\n", r)
		}
	}
}

func TestRowsParser(t *testing.T) {
	query := "SELECT id, name, active, opt_string FROM users"
	users := make([]User, 0)
	err := Parse(Query(query)).Map(usersSqlMapper(&users))
	if err != nil {
		fmt.Printf("%v", ParseErr(err))
	} else {
		fmt.Println("result:")
		for _, r := range users {
			fmt.Printf("%+v\n", r)
		}
	}
}
