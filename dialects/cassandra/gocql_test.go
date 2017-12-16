package cassandra

import (
	"fmt"
	"testing"

	"github.com/gocql/gocql"
	. "github.com/ncrypthic/dbmapper"
)

type CqlUser struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

func userCqlMapper(result *[]CqlUser) *MappedColumns {
	user := CqlUser{}
	return Columns(
		Column("id").As(&user.ID),
		Column("first_name").As(&user.FirstName),
		Column("last_name").As(&user.LastName),
		Column("email").As(&user.Email),
	).Then(func() error {
		*result = append(*result, user)
		return nil
	})
}

func userCSqlMapper(result *[]CqlUser) RowMapper {
	return func() *MappedColumns {
		return userCqlMapper(result)
	}
}

func TestCql(t *testing.T) {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "tests"
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println(err)
		panic("failed to get cassandra session ")
	}
	query := "SELECT id, first_name, last_name, email, country FROM users"
	users := make([]CqlUser, 0)
	err = Parse(session.Query(query)).Map(func() *MappedColumns {
		return userCqlMapper(&users)
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
