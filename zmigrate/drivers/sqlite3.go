package drivers

import (
	"strings"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/wyattis/z/zsql"
)

type Sqlite3Driver struct{}

func (d Sqlite3Driver) Matches(db zsql.DB) (res bool) {
	_, res = db.Driver().(*sqlite3.SQLiteDriver)
	return
}

func (d Sqlite3Driver) IsNoTableErr(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "no such table:")
}

func (d Sqlite3Driver) GetSchema(db zsql.DB) (schema Schema, err error) {
	res, err := db.Query("SELECT name, sql FROM sqlite_master WHERE type='table'")
	if err != nil {
		return
	}
	var table Table
	for res.Next() {
		if err = res.Scan(&table.Name, &table.RawSql); err != nil {
			return
		}
		schema = append(schema, table)
	}
	return
}

func init() {
	drivers = append(drivers, Sqlite3Driver{})
}
