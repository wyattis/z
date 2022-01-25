package zmigrate

import (
	"database/sql"
	"embed"
	"strings"
	"testing"

	"github.com/wyattis/z/zmigrate/drivers"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations
var migrations embed.FS

func loadSchema(path string) string {
	raw, err := migrations.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(raw), "; \n")
}
func TestToVersion1(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}
	m := New(migrations, db, nil)
	if err = m.ToVersion(1); err != nil {
		t.Error(err)
		panic(err)
	}
	schema, err := m.GetSchema()
	if err != nil {
		t.Error(err)
	}
	expected := drivers.Schema{drivers.Table{
		Name:   "user",
		RawSql: loadSchema("migrations/1_user.up.sql"),
	}}
	if !schema.Matches(expected) {
		t.Errorf("expected %v, but got %v", expected, schema)
	}
}

func TestToVersion2(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db?mode=memory")
	if err != nil {
		t.Error(err)
	}
	m := New(migrations, db, nil)
	if err = m.ToVersion(2); err != nil {
		t.Error(err)
	}
	schema, err := m.GetSchema()
	if err != nil {
		t.Error(err)
	}
	expected := drivers.Schema{drivers.Table{
		Name:   "user",
		RawSql: loadSchema("migrations/1_user.up.sql"),
	}, {
		Name:   "comment",
		RawSql: loadSchema("migrations/2_comments.up.sql"),
	}}
	if !schema.Matches(expected) {
		t.Errorf("expected %v, but got %v", expected, schema)
	}
}

func TestUpThenDown(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db?mode=memory")
	if err != nil {
		t.Error(err)
	}
	m := New(migrations, db, nil)
	if err = m.ToVersion(2); err != nil {
		t.Error(err)
	}
	if err = m.ToVersion(1); err != nil {
		t.Error(err)
	}
	schema, err := m.GetSchema()
	if err != nil {
		t.Error(err)
	}
	expected := drivers.Schema{drivers.Table{
		Name:   "user",
		RawSql: loadSchema("migrations/1_user.up.sql"),
	}}
	if !schema.Matches(expected) {
		t.Errorf("expected %v, but got %v", expected, schema)
	}
}
