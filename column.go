package dbmapper

import (
	"fmt"
)

// ColumnMap is a mapped column to target interface
type ColumnMap interface {
	// DB Column name
	Name() string
	// Maping error
	Error() error
	// Returns pointer to target
	Target() *interface{}
	// Set scan result destination
	As(target interface{}) ColumnMap
}

type column struct {
	name   string
	target interface{}
	err    error
}

func (m *column) Name() string {
	return m.name
}

func (m *column) Target() (res *interface{}) {
	if m.target != nil {
		res = &m.target
	}
	return
}

func (m *column) Error() error {
	return m.err
}

func (m *column) As(target interface{}) ColumnMap {
	if target == nil {
		m.err = fmt.Errorf("Cannot use nil to store column %s value", m.name)
	} else {
		m.target = target
	}
	return m
}

// Column returns a new default ColumnMap implementaion
func Column(name string) ColumnMap {
	return &column{name, nil, nil}
}

// Columns helper method to create slice of ColumnMap
func Columns(columns ...ColumnMap) *MappedColumns {
	return &MappedColumns{columns, func() error { return nil }}
}

// ColumnMapper provides allowing post mapping callback to process
// scan result
type ColumnMapper interface {
	Then(func()) *MappedColumns
}

// MappedColumns is a collections of ColumnMap
type MappedColumns struct {
	Columns []ColumnMap
	cb      func() error
}

// Then allows callback to proses result after row scan
func (mapped *MappedColumns) Then(cb func() error) *MappedColumns {
	mapped.cb = cb
	return mapped
}

// Done will execute callback
func (mapped *MappedColumns) Done() error {
	return mapped.cb()
}
