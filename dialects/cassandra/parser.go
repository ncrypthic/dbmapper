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

type parameter struct {
	name  string
	value interface{}
}

func (p *parameter) Name() string {
	return p.name
}

func (p *parameter) Value() (interface{}, error) {
	return p.value, nil
}

func Param(name string, val interface{}) Parameter {
	return &parameter{name, val}
}
