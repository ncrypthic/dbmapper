package cassandra

import (
	"github.com/gocql/gocql"
	. "github.com/ncrypthic/dbmapper"
)

type CqlIterator interface {
	Columns() []gocql.ColumnInfo
	Scan(...interface{}) bool
}

type CqlQuery interface {
	Iter() CqlIterator
}

type GocqlQuery interface {
	Iter() *gocql.Iter
}

// Parse is default implementation of Parser interface
func Parse(query interface{}) ResultMapper {
	return &mapper{query}
}
