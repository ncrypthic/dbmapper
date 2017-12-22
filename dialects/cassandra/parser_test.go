package cassandra

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	. "github.com/ncrypthic/dbmapper"
)

func mockRow(id uint, name string, active bool, opt_field *string) []interface{} {
	return []interface{}{id, name, active, opt_field}
}

var pos int = 0

type mockCqlIter struct {
	base *gocql.Iter
}

func resetIter() {
	pos = 0
}

func (i *mockCqlIter) Columns() []gocql.ColumnInfo {
	return []gocql.ColumnInfo{
		gocql.ColumnInfo{Name: "id"},
		gocql.ColumnInfo{Name: "name"},
		gocql.ColumnInfo{Name: "active"},
		gocql.ColumnInfo{Name: "opt_field"},
	}
}

func (i *mockCqlIter) Scan(result ...interface{}) bool {
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
	if pos >= len(rows) {
		return false
	} else {
		for j, v := range rows[pos] {
			i.assignVal(result[j], v)
		}
		pos++
		return true
	}
}

func (i *mockCqlIter) assignVal(dest interface{}, source interface{}) {
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

type mockCqlQuery struct{}

func (i *mockCqlQuery) Iter() CqlIterator {
	return &mockCqlIter{nil}
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
	return &mockCqlQuery{}
}

func TestRowParser(t *testing.T) {
	defer resetIter()
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
	defer resetIter()
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

func TestQueryMapper(t *testing.T) {
	namedSql := "insert into test values (:id, :name)"
	q := Prepare(namedSql).With(
		Parameter{"id", "123"},
		Parameter{"name", 1},
		Parameter{"phone", "0827126"},
	)
	expectedParams := []interface{}{"123", 1}
	expectedParamNames := []string{":id", ":name"}
	if len(q.ParamNames()) != 2 {
		t.Errorf("Fail: expect %v parameters, got %v instead", expectedParamNames, q.ParamNames())
	}
	if len(q.Params()) != 2 {
		t.Errorf("Fail: expect %v parameters, got %v instead", expectedParams, q.Params())
	}
}
