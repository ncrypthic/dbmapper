package dbmapper

import (
	"database/sql"
	"errors"
	"fmt"
)

// NoResultErr no result from query
type NoResultErr error

// RowMapper is a wrapper to create target struct and
// returns slice of ColumnMap for columns scan
type RowMapper func() MappedColumns

// Place holder for unmapped column
type dummy struct{}

func (d *dummy) Scan(_ interface{}) error {
	return nil
}

type mapper struct {
	rows *sql.Rows
	err  error
}

func (m *mapper) targets(mapTarget map[string]*interface{}, names []string) []interface{} {
	result := make([]interface{}, len(names))
	for i, name := range names {
		if target, ok := mapTarget[name]; ok {
			result[i] = *target
		} else {
			result[i] = new(dummy)
		}
	}
	return result
}

func (m *mapper) Map(rowMapper RowMapper) (mapErr error) {
	if m.err != nil {
		fmt.Printf("%+v\n", m.err)
		return m.err
	}
	var dbColumns []string
	rowMap := rowMapper()
	if m.rows == nil {
		return NoResultErr(errors.New("No rows from query"))
	}
	for m.rows.Next() {
		targets := rowMap.columns
		targetMap := make(map[string]*interface{})
		for _, column := range targets {
			if columnErr := column.Error(); columnErr != nil {
				mapErr = columnErr
				fmt.Printf("%+v\n", mapErr)
			}
			targetMap[column.Name()] = column.Target()
		}
		if mapErr != nil {
			fmt.Printf("%+v\n", mapErr)
			return mapErr
		}
		if dbColumns == nil {
			dbColumns, mapErr = m.rows.Columns()
		}
		if mapErr != nil {
			fmt.Printf("%+v\n", mapErr)
			return
		}
		dest := m.targets(targetMap, dbColumns)
		if mapErr = m.rows.Scan(dest...); mapErr != nil {
			fmt.Printf("%+v\n", mapErr)
			return
		}
		if mapErr = rowMap.done(); mapErr != nil {
			return
		}
	}
	return
}
