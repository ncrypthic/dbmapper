package mysql

import (
	"database/sql"

	. "github.com/ncrypthic/dbmapper"
)

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
		return m.err
	}
	var dbColumns []string
	rowMap := rowMapper()
	for m.rows.Next() {
		targets := rowMap.Columns
		targetMap := make(map[string]*interface{})
	TargetLoop:
		for _, column := range targets {
			if columnErr := column.Error(); columnErr != nil {
				mapErr = columnErr
				break TargetLoop
			}
			targetMap[column.Name()] = column.Target()
		}
		if mapErr != nil {
			return mapErr
		}
		if dbColumns == nil {
			dbColumns, mapErr = m.rows.Columns()
		}
		if mapErr != nil {
			return
		}
		dest := m.targets(targetMap, dbColumns)
		mapErr = m.rows.Scan(dest...)
		if mapErr != nil {
			return
		}
		mapErr = rowMap.Done()
		if mapErr != nil {
			return
		}
	}
	return
}
