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
	idRows := sqlMock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3)
	mock.ExpectQuery("SELECT id FROM users").WillReturnRows(idRows)
	res := make([]int32, 0)
	err = Parse(db.Query("SELECT id FROM users")).Map(Int32("id", &res))
	if err != nil {
		fmt.Printf("%v", ParseErr(err))
	} else {
		fmt.Println("result:")
		for _, r := range res {
			fmt.Printf("%+v\n", r)
		}
	}
	nameRows := sqlMock.NewRows([]string{"name"}).AddRow("alice").AddRow("bow").AddRow("charlie")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(nameRows)
	names := make([]string, 0)
	err = Parse(db.Query("SELECT name FROM users")).Map(String("name", &names))
	if err != nil {
		fmt.Printf("%v", ParseErr(err))
	} else {
		fmt.Println("result:")
		for _, r := range names {
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
	namedSql := "insert into test(id, name, created) values (:id, :name, NOW())"
	q := Prepare(namedSql).With(
		Param("id", "123"),
		Param("name", 1),
		Param("phone", "0827126"),
	)
	expectedParams := []interface{}{"123", 1}
	expectedParamNames := []string{":id", ":name"}
	expectedSql := "insert into test(id, name, created) values (?, ?, NOW())"
	if len(q.ParamNames()) != 2 {
		t.Errorf("Fail: expect %v parameters, got %v instead", expectedParamNames, q.ParamNames())
	}
	if len(q.Params()) != 2 {
		t.Errorf("Fail: expect %v parameters, got %v instead", expectedParams, q.Params())
	}
	if q.SQL() != expectedSql {
		t.Errorf("Fail: expect [ %v ] sql string, got [ %v ] instead", expectedSql, q.SQL())
	}
	ids := []interface{}{"1", "2"}
	selectSql := "select id, name, phone from test where id IN (:ids) and name like :keyword"
	q = Prepare(selectSql).With(
		Param("ids", ids...),
		Param("keyword", "%abc%"),
	)
	expectedParams = []interface{}{"1", "2", "%abc%"}
	expectedParamNames = []string{":ids", ":keyword"}
	expectedSql = "select id, name, phone from test where id IN (?, ?) and name like ?"
	if len(q.ParamNames()) != 2 {
		t.Errorf("Fail: expect %v parameters, got %v instead", expectedParamNames, q.ParamNames())
	}
	if len(q.Params()) != 3 {
		t.Errorf("Fail: expect %v parameters, got %v instead", expectedParams, q.Params())
	}
	if q.SQL() != expectedSql {
		t.Errorf("Fail: expect [ %v ] sql string, got [ %v ] instead", expectedSql, q.SQL())
	}
	for i := 0; i < 10; i++ {
		for idx, queryParam := range q.Params() {
			if expectedParams[idx] != queryParam {
				t.Errorf("Fail: expect [ %v ] sql string, got [ %v ] instead", expectedParams, q.Params())
			}
		}
	}
}
