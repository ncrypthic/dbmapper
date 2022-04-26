package cassandra

import (
	"errors"

	"github.com/gocql/gocql"
	. "github.com/ncrypthic/dbmapper"
)

var (
	ErrNoRows error = errors.New("no results")
)

type CqlIterator interface {
	Columns() []gocql.ColumnInfo
	Scan(...interface{}) bool
	NumRows() int
	Close() error
}

type CqlQuery interface {
	Iter() CqlIterator
}

type cqlQuery struct {
	query *gocql.Query
}

func (q *cqlQuery) Iter() CqlIterator {
	return q.query.Iter()
}

// Parse is default implementation of Parser interface
func Parse(q *gocql.Query) ResultMapper {
	query := &cqlQuery{query: q}
	return &mapper{query}
}

// Parse is default implementation of Parser interface
func ParseCqlQuery(q CqlQuery) ResultMapper {
	return &mapper{query: q}
}
