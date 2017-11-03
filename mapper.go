package sqlmapper

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
