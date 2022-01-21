package zmigrate

import (
	"crypto/md5"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wyattis/z/zmigrate/drivers"
	"github.com/wyattis/z/zsql"
	"github.com/wyattis/z/zstring"
)

//go:embed schema.sql
var schema string

var (
	ErrInvalidVersion   = errors.New("invalid version number")
	ErrMigrationChanged = errors.New("a migration has changed since it was originally run")
)

type Config struct {
	Table           string
	SkipTransaction bool
}

func (c *Config) applyDefaults() {
	if c.Table == "" {
		c.Table = "migrations"
	}
	// if c.DirName == "" {
	// 	c.DirName = "migrations"
	// }
}

type migration struct {
	Id        int
	File      string
	Name      string
	SQL       string
	Md5       string
	CreatedAt uint64 `db:"created_at"`
}

func New(source fs.ReadDirFS, db *sql.DB, config *Config) *Migrator {
	if config == nil {
		config = &Config{}
	}
	config.applyDefaults()
	return &Migrator{Source: source, db: db, config: *config, mut: &sync.Mutex{}}
}

type Migrator struct {
	Source        fs.ReadDirFS
	db            *sql.DB
	config        Config
	isInitialized bool
	mut           *sync.Mutex
	Driver        drivers.Driver

	currentId int
}

type Exec interface {
	Exec(q string, arguments ...interface{}) (sql.Result, error)
}

func (m *Migrator) GetSchema() (schema drivers.Schema, err error) {
	if err = m.init(); err != nil {
		return
	}
	res, err := m.Driver.GetSchema(m.db)
	if err != nil {
		return
	}
	for _, table := range res {
		if table.Name != m.config.Table {
			schema = append(schema, table)
		}
	}
	return
}

// Migrate to the latest version
func (m *Migrator) ToLatest() (err error) {
	if err = m.init(); err != nil {
		return
	}
	// TODO: determine the latest version id and call `m.ToVersion`
	return
}

// Migrate up or down to a specific state in the db. This allows you to work on
// your schema without ruining existing environments by locking it to a version.
func (m *Migrator) ToVersion(version int) (err error) {
	if err = m.init(); err != nil {
		return
	}
	currentVersion, err := m.currentVersion()
	if err != nil || version == currentVersion {
		return
	}
	availableUp, availableDown, err := m.getAvailable()
	if err != nil {
		return
	}
	// determine if this version exists. Throw error if it doesn't.
	if !m.versionExists(version, availableUp) {
		return ErrInvalidVersion
	}

	existingUp, err := m.getExisting()
	if err != nil {
		return
	}
	// compare all md5 hashes for all 'up' migrations with the existing ones.
	if !m.hashesMatch(availableUp[:currentVersion], existingUp[:currentVersion]) {
		return ErrMigrationChanged
	}
	// if there are conflicts we throw an error.
	if version > currentVersion {
		return m.withTx(func(tx zsql.Exec) error {
			return m.up(tx, availableUp[currentVersion:version])
		})
	} else {
		return m.withTx(func(tx zsql.Exec) error {
			return m.down(tx, m.reverse(availableDown[version:currentVersion]))
		})
	}
}

func (m *Migrator) withTx(handler zsql.TxHandler) (err error) {
	if m.config.SkipTransaction {
		return zsql.WithoutTx(m.db, handler)
	} else {
		return zsql.WithTx(m.db, handler)
	}
}

func (m *Migrator) up(tx Exec, migrations []migration) (err error) {
	if err != nil {
		return
	}
	qInsert := fmt.Sprintf("INSERT INTO %s (file, name, md5, created_at) VALUES (?, ?, ?, ?)", m.config.Table)
	for _, m := range migrations {
		if _, err = tx.Exec(m.SQL); err != nil {
			return
		}
		if _, err = tx.Exec(qInsert, m.File, m.Name, m.Md5, time.Now()); err != nil {
			return
		}
	}
	return
}

func (m *Migrator) down(tx Exec, migrations []migration) (err error) {
	qDelete := fmt.Sprintf("DELETE FROM %s where id=?", m.config.Table)
	for _, m := range migrations {
		if _, err = tx.Exec(m.SQL); err != nil {
			return
		}
		if _, err = tx.Exec(qDelete, m.Id); err != nil {
			return
		}
	}
	return
}

func (m *Migrator) hashesMatch(a []migration, b []migration) bool {
	for i := range a {
		if strings.TrimSpace(a[i].Md5) != strings.TrimSpace(b[i].Md5) {
			return false
		}
	}
	return true
}

func (m *Migrator) currentVersion() (res int, err error) {
	existing, err := m.getExisting()
	if err == nil && len(existing) > 0 {
		res = existing[len(existing)-1].Id
	}
	return
}

func (m *Migrator) reverse(migrations []migration) (result []migration) {
	for i := range migrations {
		result = append(result, migrations[len(migrations)-i-1])
	}
	return
}

func (m *Migrator) versionExists(version int, available []migration) bool {
	versionExists := false
	for _, mig := range available {
		versionExists = mig.Id == version
		if versionExists {
			break
		}
	}
	return versionExists
}

// initialize the migrations table
func (m *Migrator) init() (err error) {
	m.mut.Lock()
	defer m.mut.Unlock()
	if m.isInitialized {
		return
	}
	// determine the SQL driver flavor
	if m.Driver, err = drivers.Get(m.db); err != nil {
		return
	}

	// create migration table if it doesn't already exist
	_, err = m.getExisting()
	if m.Driver.IsNoTableErr(err) {
		schema = fmt.Sprintf(schema, m.config.Table)
		if _, err = m.db.Exec(schema); err != nil {
			return
		}
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	m.isInitialized = true
	return
}

// get all existing migrations
func (m *Migrator) getExisting() (res []migration, err error) {
	q := fmt.Sprintf("SELECT id, file, name, md5, created_at FROM %s ORDER BY id ASC", m.config.Table)
	rows, err := m.db.Query(q)
	if err != nil {
		return
	}
	for rows.Next() {
		if rows.Err() != nil {
			err = rows.Err()
			return
		}
		mig := migration{}
		rows.Scan(&mig.Id, &mig.File, &mig.Name, &mig.Md5, &mig.CreatedAt)
		res = append(res, mig)
	}
	return
}

func (m *Migrator) getAvailable() (up []migration, down []migration, err error) {
	entries, err := m.Source.ReadDir(".")
	if err != nil {
		return
	}
	if len(entries) == 1 && entries[0].IsDir() {
		source, err := fs.Sub(m.Source, entries[0].Name())
		if err != nil {
			return up, down, err
		}
		m.Source = source.(fs.ReadDirFS)
		entries, err = m.Source.ReadDir(".")
	}
	if err != nil {
		return
	}
	for _, dir := range entries {
		if dir.IsDir() {
			continue
		}
		version, name, found := zstring.Cut(dir.Name(), "_")
		if !found {
			err = fmt.Errorf("file name should have format 1_name.up.sql or 1_name.down.sql instead of %s", dir.Name())
			return
		}
		id, err2 := strconv.Atoi(version)
		if err2 != nil {
			err = fmt.Errorf("version should be a valid integer instead of %s", version)
			return
		}
		f, err2 := m.Source.Open(dir.Name())
		if err2 != nil {
			err = err2
			return
		}
		defer f.Close()
		content, err2 := io.ReadAll(f)
		if err2 != nil {
			err = err2
			return
		}
		hash := md5.Sum(content)
		migration := migration{
			Id:   id,
			Name: name,
			File: dir.Name(),
			Md5:  string(hash[:]),
			SQL:  string(content),
		}
		if strings.HasSuffix(name, ".up.sql") {
			up = append(up, migration)
		} else if strings.HasSuffix(name, ".down.sql") {
			down = append(down, migration)
		} else {
			err = fmt.Errorf("invalid filename %s", name)
			return
		}
	}
	return
}
