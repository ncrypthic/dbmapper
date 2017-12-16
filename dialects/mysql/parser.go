package mysql

import (
	"database/sql"

	. "github.com/ncrypthic/dbmapper"
)

// ResultMapper is database result set mapper
type ResultMapper interface {
	Map(RowMapper) error
}

// Parse is default implementation of Parser interface
func Parse(rows *sql.Rows, err error) ResultMapper {
	return &mapper{rows, err}
}
