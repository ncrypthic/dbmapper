package mysql

import (
	"database/sql"
	"regexp"

	. "github.com/ncrypthic/dbmapper"
)

// Parse is default implementation of Parser interface
func Parse(rows *sql.Rows, err error) ResultMapper {
	return &mapper{rows, err}
}
