package dbmapper

import (
	"regexp"
)

// NoResultErr no result from query
type NoResultErr error

// MapScanErr  no result from query
type MapScanErr error

type ResultSet interface {
	Next()
	Columns() []string
	Scan(...interface{}) error
}

// RowMapper is a wrapper to create target struct and
// returns slice of ColumnMap for columns scan
type RowMapper func() *MappedColumns

// Parameter is query parameter key value pair
type Parameter struct {
	Name  string
	Value interface{}
}

type QueryMapper interface {
	With(...Parameter) QueryMapper
	Params() []interface{}
	SQL() string
	ParamNames() []string
}

// ResultMapper is database result set mapper
type ResultMapper interface {
	Map(RowMapper) error
}

type query struct {
	namedSql   string
	sql        string
	params     map[string]interface{}
	paramNames []string
}

func (q *query) Params() []interface{} {
	args := make([]interface{}, 0)
	for _, name := range q.paramNames {
		if val, ok := q.params[name]; ok {
			args = append(args, val)
		}
	}
	return args
}

func (q *query) SQL() string {
	return q.sql
}

func (q *query) ParamNames() []string {
	return q.paramNames
}

func (q *query) With(parameters ...Parameter) QueryMapper {
	pattern := regexp.MustCompile(":([a-z0-9-_]+)")
	q.sql = pattern.ReplaceAllString(q.sql, "?")
	if q.params == nil {
		q.params = make(map[string]interface{})
	}
	for _, p := range parameters {
		q.params[":"+p.Name] = p.Value
	}
	return q
}

// Prepare
func Prepare(namedSql string) QueryMapper {
	pattern := regexp.MustCompile(":([a-z0-9-_]+)")
	paramNames := pattern.FindAllString(namedSql, -1)
	sql := pattern.ReplaceAllString(namedSql, "?")
	return &query{namedSql: namedSql, sql: sql, paramNames: paramNames}
}
