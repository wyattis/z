package drivers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/wyattis/z/zsql"
)

type PostgresDriver struct{}

func (d PostgresDriver) Matches(db zsql.DB) (res bool) {
	_, res = db.Driver().(*pq.Driver)
	return
}

func (d PostgresDriver) IsNoTableErr(err error) bool {
	return err != nil && strings.HasSuffix(err.Error(), "does not exist")
}

func (d PostgresDriver) GetSchema(db zsql.DB) (schema Schema, err error) {
	tables := []string{}
	res, err := db.Query("select table_schema from information_schema.tables")
	if err != nil {
		return
	}
	for res.Next() {
		var s string
		if err = res.Scan(&s); err != nil {
			return
		}
		tables = append(tables, s)
	}
	// TODO
	return
}

func (d PostgresDriver) ExpandError(err error) error {
	if pqerr, ok := err.(pq.Error) {
		return errors.New(fmt.Sprintln(pqerr.Line, pqerr.Position, pqerr.Message))
	} else if pqerr , ok := err.(pq.Err) {
		
	}
	return err
}

func init() {
	drivers = append(drivers, PostgresDriver{})
}
