package cassandra

import (
	"errors"
	"fmt"

	. "github.com/ncrypthic/sqlmapper"
)

// Place holder for unmapped column
type dummy struct{}

func (d *dummy) Scan(_ interface{}) error {
	return nil
}

type mapper struct {
	query CqlQuery
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
	rowMap := rowMapper()
	if m.query == nil {
		return NoResultErr(errors.New("Query is not valid"))
	}
	var dbColumns []string
	rs := m.query.Iter()
	for {
		if dbColumns == nil {
			dbColumns = make([]string, 0)
			for _, cqlColumn := range rs.Columns() {
				dbColumns = append(dbColumns, cqlColumn.Name)
			}
		}
		targets := rowMap.Columns
		targetMap := make(map[string]*interface{})
		for _, column := range targets {
			if columnErr := column.Error(); columnErr != nil {
				mapErr = columnErr
				fmt.Printf("%+v\n", mapErr)
			}
			targetMap[column.Name()] = column.Target()
		}
		dest := m.targets(targetMap, dbColumns)
		if scanOk := rs.Scan(dest...); !scanOk {
			break
		}
		if mapErr = rowMap.Done(); mapErr != nil {
			return mapErr
		}
	}
	return
}
