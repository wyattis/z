package drivers

import (
	"database/sql"
	"errors"
)

type Driver interface {
	Matches(*sql.DB) bool
	GetSchema(*sql.DB) (Schema, error)
	IsNoTableErr(error) bool
}

type Schema []Table

func (s Schema) Matches(other Schema) bool {
	if len(s) != len(other) {
		return false
	}
	for i := range s {
		if !s[i].Matches(other[i]) {
			return false
		}
	}
	return true
}

type Table struct {
	Name    string
	RawSql  string
	Columns []Column
	Indices []Index
}

func (t *Table) Matches(other Table) bool {
	if t.Name != other.Name || t.RawSql != other.RawSql || len(t.Columns) != len(other.Columns) || len(t.Indices) != len(other.Indices) {
		return false
	}
	for i := range t.Columns {
		if t.Columns[i] != other.Columns[i] {
			return false
		}
	}
	for i := range t.Indices {
		if !t.Indices[i].Matches(other.Indices[i]) {
			return false
		}
	}
	return true
}

type Column struct {
	Name     string
	Type     string
	Nullable bool
}

type Index struct {
	Name    string
	Type    string
	Columns []string
}

func (i *Index) Matches(other Index) bool {
	if i.Name != other.Name || i.Type != other.Type {
		return false
	}
	for j := range i.Columns {
		if i.Columns[j] != other.Columns[j] {
			return false
		}
	}
	return true
}

var drivers = []Driver{}

func Get(db *sql.DB) (Driver, error) {
	for _, d := range drivers {
		if d.Matches(db) {
			return d, nil
		}
	}
	return nil, errors.New("no matching driver found")
}
