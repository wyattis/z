package zmigrate

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/wyattis/z/zsql"
)

func TestSeedSingle(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}
	seeder, err := NewSeeder(db, nil)
	if err != nil {
		t.Error(err)
	}
	seeder.AddSeed("test", 1, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		if _, err = tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS word  (id INTEGER PRIMARY KEY, word TEXT);"); err != nil {
			return
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?)", "hello", "world")
		return
	})

	if err = seeder.SeedTo("test", 1); err != nil {
		t.Error(err)
	}
}

func TestSeedUpDown(t *testing.T) {
	db, seeder, err := makeSeeder()
	if err != nil {
		t.Error(err)
	}
	seeder.AddSeed("test", 2, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?), (?)", "I", "am", "here")
		return
	})
	seeder.AddSeed("test", 1, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?)", "hello", "world")
		return
	})
	if err = seeder.SeedTo("test", 1); err != nil {
		t.Error(err)
	}
	expected := []string{"hello", "world"}
	if err = TableColumnMatches(db, "word", "word", expected); err != nil {
		t.Error(err)
	}
	if err = seeder.SeedTo("test", 2); err != nil {
		t.Error(err)
	}

	expected = []string{"hello", "world", "I", "am", "here"}
	if err = TableColumnMatches(db, "word", "word", expected); err != nil {
		t.Error(err)
	}

	// should fail to decrease the seed version
	if err = seeder.SeedTo("test", 1); err == nil {
		t.Error("should throw an error when we attempt to decrease the seed")
	}
}

func TestSeedDouble(t *testing.T) {
	db, seeder, err := makeSeeder()
	if err != nil {
		t.Error(err)
	}
	seeder.AddSeed("test", 2, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?), (?)", "I", "am", "here")
		return
	})
	seeder.AddSeed("test", 1, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?)", "hello", "world")
		return
	})
	if err = seeder.SeedTo("test", 2); err != nil {
		t.Error(err)
	}

	expected := []string{"hello", "world", "I", "am", "here"}
	if err = TableColumnMatches(db, "word", "word", expected); err != nil {
		t.Error(err)
	}
}

func makeSeeder() (db *sql.DB, seeder *Seeder, err error) {
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS word  (id INTEGER PRIMARY KEY, word TEXT);"); err != nil {
		return
	}
	seeder, err = NewSeeder(db, nil)
	if err != nil {
		return
	}
	return
}

func TableColumnMatches(db zsql.DB, table string, column string, expected []string) (err error) {
	rows, err := SelectAll(db, table, column)
	if err != nil {
		return
	}
	if !reflect.DeepEqual(rows, expected) {
		err = fmt.Errorf("Expected %s, but got %s", expected, rows)
	}
	return
}

func SelectAll(db zsql.DB, table string, column string) (rows []string, err error) {
	res, err := db.Query(fmt.Sprintf("SELECT %s FROM %s", column, table))
	if err != nil {
		return
	}
	for res.Next() {
		row := ""
		if err = res.Scan(&row); err != nil {
			return
		}
		rows = append(rows, row)
	}
	return
}
