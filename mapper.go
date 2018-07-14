package dbmapper

import (
	"errors"
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

// Int32 returns a row mapper for single int32 column
func String(columnName string, dst *[]string) RowMapper {
	data := ""
	return func() *MappedColumns {
		return Columns(
			Column(columnName).As(&data),
		).Then(func() error {
			*dst = append(*dst, data)
			return nil
		})
	}
}

// Int32 returns a row mapper for single int32 column
func Int32(columnName string, dst *[]int32) RowMapper {
	data := int32(0)
	return func() *MappedColumns {
		return Columns(
			Column(columnName).As(&data),
		).Then(func() error {
			*dst = append(*dst, data)
			return nil
		})
	}
}

// Int64 returns a row mapper for single int32 column
func Int64(columnName string, dst *[]int64) RowMapper {
	data := int64(0)
	return func() *MappedColumns {
		return Columns(
			Column(columnName).As(&data),
		).Then(func() error {
			*dst = append(*dst, data)
			return nil
		})
	}
}

type query struct {
	namedSql    string
	sql         string
	params      map[string]interface{}
	paramValues []interface{}
	paramNames  []string
	err         error
}

func (q *query) Error() error {
	return q.err
}

func (q *query) Params() []interface{} {
	return q.paramValues
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

func (q *query) getParameter(name string, params []Parameter) ([]interface{}, error) {
	for _, p := range params {
		if name == p.Name() {
			return p.Value()
		}
	}
	return nil, nil
}

func (q *query) With(parameters ...Parameter) QueryMapper {
	q.sql = q.namedSql
	if q.params == nil {
		q.params = make(map[string]interface{})
	}
	pattern, err := regexp.Compile(":([a-zA-Z0-9_]+)")
	if err != nil {
		q.err = err
		return q
	}
	res := pattern.FindAllStringSubmatch(q.namedSql, -1)
	for _, match := range res {
		if len(match) == 2 {
			paramName := match[1]
			value, err := q.getParameter(paramName, parameters)
			if err != nil {
				q.err = err
				return q
			}
			if len(value) > 1 {
				q.paramNames = append(q.paramNames, paramName)
				sliceElmts := []string{}
				for idx, elmt := range value {
					q.params[paramName+"_"+string(idx)] = elmt
					q.paramValues = append(q.paramValues, elmt)
					sliceElmts = append(sliceElmts, "?")
				}
				q.sql = strings.Replace(q.sql, ":"+paramName, strings.Join(sliceElmts, ", "), 1)
			} else if len(value) == 1 {
				q.paramNames = append(q.paramNames, paramName)
				q.params[paramName] = value[0]
				q.paramValues = append(q.paramValues, value[0])
				q.sql = strings.Replace(q.sql, ":"+paramName, "?", 1)
			} else {
				q.err = errors.New("Missing paramters")
				return q
			}
		} else {
			continue
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
