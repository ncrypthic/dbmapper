package dbmapper

import (
	"fmt"
	"testing"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
)

type ParseErr error

type User struct {
	ID   string
	Name string
}

type UsersParser struct {
	result []User
}

func (p *UsersParser) Result() []User {
	return p.result
}

func (p *UsersParser) parseRow() MappedColumns {
	row := User{}
	return Columns(
		Column("id").As(&row.ID),
		Column("name").As(&row.Name),
	).Then(func() error {
		p.result = append(p.result, row)
		return nil
	})
}

func (p *UsersParser) Parse() RowMapper {
	p.result = make([]User, 0)
	return p.parseRow
}

func TestParser(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		panic("failed to mock database")
	}
	defer db.Close()
	rows := sqlMock.NewRows([]string{"id", "name"}).
		AddRow(1, "alice").
		AddRow(2, "bob").
		AddRow(3, "charlie")
	mock.ExpectQuery("SELECT id, name FROM users").WillReturnRows(rows)
	usersParser := UsersParser{}
	if err := Parse(db.Query("SELECT id, name FROM users")).Map(usersParser.Parse()); err != nil {
		fmt.Printf("%v", ParseErr(err))
	} else {
		fmt.Println("result:")
		for _, r := range usersParser.Result() {
			fmt.Printf("%+v\n", r)
		}
	}
}
