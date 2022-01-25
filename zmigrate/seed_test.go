package zmigrate

import (
	"context"
	"database/sql"
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

func TestSeedDouble(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS word  (id INTEGER PRIMARY KEY, word TEXT);"); err != nil {
		t.Error(err)
	}
	seeder, err := NewSeeder(db, nil)
	if err != nil {
		t.Error(err)
	}
	seeder.AddSeed("test", 1, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?)", "hello", "world")
		return
	})
	seeder.AddSeed("test", 2, func(ctx context.Context, tx zsql.Tx, seed Seed) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO word (word) VALUES (?), (?), (?)", "I", "am", "here")
		return
	})

	if err = seeder.SeedTo("test", 1); err != nil {
		t.Error(err)
	}
	if err = seeder.SeedTo("test", 2); err != nil {
		t.Error(err)
	}
	// should fail to decrease the seed version
	if err = seeder.SeedTo("test", 1); err == nil {
		t.Error("should throw an error when we attempt to decrease the seed")
	}
}
