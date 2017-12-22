package mysql

import (
	"database/sql"

	. "github.com/ncrypthic/dbmapper"
)

// Parse is default implementation of Parser interface
func Parse(rows *sql.Rows, err error) ResultMapper {
	return &mapper{rows, err}
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
