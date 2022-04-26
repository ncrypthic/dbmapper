package cassandra

import (
	"fmt"

	. "github.com/ncrypthic/dbmapper"
)

type mapper struct {
	query CqlQuery
}

func (m *mapper) targets(mapTarget map[string]*interface{}, names []string) []interface{} {
	result := make([]interface{}, len(names))
	for i, name := range names {
		if target, ok := mapTarget[name]; ok {
			result[i] = *target
		} else {
			result[i] = nil
		}
	}
	return result
}

func (m *mapper) Map(rowMapper RowMapper) (mapErr error) {
	rowMap := rowMapper()
	var dbColumns []string
	rs := m.query.Iter()
	if rs.NumRows() == 0 {
		return ErrNoRows
	}
	for {
		dbColumns = make([]string, 0)
		for _, cqlColumn := range rs.Columns() {
			dbColumns = append(dbColumns, cqlColumn.Name)
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
	return rs.Close()
}
