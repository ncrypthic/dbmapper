package mysql

import (
	"database/sql"

	. "github.com/ncrypthic/sqlmapper"
)

// ResultMapper is database result set mapper
type ResultMapper interface {
	Map(RowMapper) error
}

// Parser returns a RowMapper object
type Parser func(*sql.Rows, error) RowMapper

// Parse is default implementation of Parser interface
func Parse(rows *sql.Rows, err error) ResultMapper {
	return &mapper{rows, err}
}
