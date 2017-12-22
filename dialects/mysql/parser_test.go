package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/ncrypthic/dbmapper"
)

type ParseErr error

type User struct {
	ID        string
	Name      string
	Active    bool
	OptString sql.NullString
}

func userSqlMapper(result *[]User) *MappedColumns {
	user := User{}
	return Columns(
		Column("id").As(&user.ID),
		Column("name").As(&user.Name),
		Column("active").As(&user.Active),
		Column("opt_string").As(&user.OptString),
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

func TestRowParser(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		panic("failed to mock database")
	}
	defer db.Close()
	rows := sqlMock.NewRows([]string{"id", "name", "active", "opt_string"}).
		AddRow(1, "alice", true, nil).
		AddRow(2, "bob", true, "11111111").
		AddRow(3, "charlie", false, nil)
	mock.ExpectQuery("SELECT id, name, active, opt_string FROM users").WillReturnRows(rows)
	users := make([]User, 0)
	err = Parse(db.Query("SELECT id, name, active, opt_string FROM users")).Map(func() *MappedColumns {
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
	db, mock, err := sqlMock.New()
	if err != nil {
		panic("failed to mock database")
	}
	defer db.Close()
	rows := sqlMock.NewRows([]string{"id", "name", "active", "opt_string"}).
		AddRow(1, "alice", true, nil).
		AddRow(2, "bob", true, "11111111").
		AddRow(3, "charlie", false, nil)
	mock.ExpectQuery("SELECT id, name, active, opt_string FROM users").WillReturnRows(rows)
	users := make([]User, 0)
	err = Parse(db.Query("SELECT id, name, active, opt_string FROM users")).Map(usersSqlMapper(&users))
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
		Param("id", "123"),
		Param("name", 1),
		Param("phone", "0827126"),
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
