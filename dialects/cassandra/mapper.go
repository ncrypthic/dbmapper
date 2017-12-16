package cassandra

import (
	"errors"
	"fmt"

	. "github.com/ncrypthic/dbmapper"
)

type mapper struct {
	query interface{}
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
	if m.query == nil {
		return NoResultErr(errors.New("Query is not valid"))
	}
	var dbColumns []string
	var rs CqlIterator
	switch t := m.query.(type) {
	case CqlQuery:
		rs = t.Iter()
	case GocqlQuery:
		rs = t.Iter()
	default:
		return fmt.Errorf("Failed to iterate query result")
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
	return
}
