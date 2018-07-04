package dbmapper

import (
	"log"
	"regexp"
	"strings"
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

type Parameter interface {
	Name() string
	Value() ([]interface{}, error)
}

type QueryMapper interface {
	With(...Parameter) QueryMapper
	Params() []interface{}
	SQL() string
	ParamNames() []string
	Error() error
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
	err        error
}

func (q *query) Error() error {
	return q.err
}

func (q *query) Params() []interface{} {
	params := make([]interface{}, 0)
	if q.err != nil {
		log.Printf("Failed to get query string, %v", q.err)
		return params
	}
	for _, val := range q.params {
		params = append(params, val)
	}
	return params
}

func (q *query) SQL() string {
	if q.err != nil {
		log.Printf("Failed to get query string, %v", q.err)
		return ""
	}
	return q.sql
}

func (q *query) ParamNames() []string {
	if q.err != nil {
		log.Printf("Failed to get query string, %v", q.err)
		return make([]string, 0)
	}
	return q.paramNames
}

func (q *query) With(parameters ...Parameter) QueryMapper {
	q.sql = q.namedSql
	if q.params == nil {
		q.params = make(map[string]interface{})
	}
	for _, p := range parameters {
		paramName := ":" + p.Name()
		if strings.Index(q.namedSql, paramName) == -1 {
			continue
		}
		q.paramNames = append(q.paramNames, paramName)
		val, err := p.Value()
		if err != nil {
			q.err = err
			return q
		}
		if len(val) > 1 {
			sliceElmts := []string{}
			for idx, elmt := range val {
				q.params[paramName+"_"+string(idx)] = elmt
				sliceElmts = append(sliceElmts, "?")
			}
			q.sql = strings.Replace(q.sql, paramName, strings.Join(sliceElmts, ", "), 1)
		} else if len(val) == 1 {
			q.params[paramName] = val
			q.sql = strings.Replace(q.sql, paramName, "?", 1)
		}
	}
	return q
}

// Prepare
func Prepare(namedSql string) QueryMapper {
	pattern := regexp.MustCompile(":([a-z0-9-_]+)")
	sql := pattern.ReplaceAllString(namedSql, "?")
	return &query{namedSql: namedSql, sql: sql, paramNames: make([]string, 0)}
}

type parameter struct {
	name  string
	value []interface{}
}

func (p *parameter) Name() string {
	return p.name
}

func (p *parameter) Value() ([]interface{}, error) {
	return p.value, nil
}

func Param(name string, val ...interface{}) Parameter {
	return &parameter{name, val}
}
