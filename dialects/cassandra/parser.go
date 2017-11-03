package cassandra

import (
	"github.com/gocql/gocql"
	. "github.com/ncrypthic/sqlmapper"
)

type CqlIterator interface {
	Columns() []gocql.ColumnInfo
	Scan(...interface{}) bool
}

type CqlQuery interface {
	Iter() CqlIterator
}

// ResultMapper is database result set mapper
type ResultMapper interface {
	Map(RowMapper) error
}

// Parse is default implementation of Parser interface
func Parse(query CqlQuery) ResultMapper {
	return &mapper{query}
}
