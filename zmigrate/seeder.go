package zmigrate

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/wyattis/z/zmigrate/drivers"
	"github.com/wyattis/z/zsql"
)

//go:embed seed-schema.sql
var seedSchema string

type SeederConfig struct {
	TableName string
}

func (s *SeederConfig) validate() error {
	if s.TableName == "" {
		s.TableName = "seeds"
	}
	return nil
}

func NewSeeder(db zsql.DB, config *SeederConfig) (s *Seeder, err error) {
	if config == nil {
		config = &SeederConfig{}
	}
	if err = config.validate(); err != nil {
		return
	}
	s = &Seeder{db: db, mut: &sync.Mutex{}, config: *config, seeds: make(map[string][]Seed)}
	s.driver, err = drivers.Get(db)
	return
}

type SeedHandler func(ctx context.Context, tx zsql.Tx, seed Seed) error
type SeedHandlerx func(ctx context.Context, tx zsql.Txx, seed Seed) error

type Seed struct {
	Name    string
	Version int
	Handler SeedHandler
}

type Seeder struct {
	db            zsql.DB
	config        SeederConfig
	driver        drivers.Driver
	isInitialized bool
	mut           *sync.Mutex
	seeds         map[string][]Seed
}

// Add a seed handler for the specified name/version
func (s *Seeder) AddSeed(name string, version int, handler SeedHandler) *Seeder {
	if _, exists := s.seeds[name]; !exists {
		s.seeds[name] = make([]Seed, 0)
	}
	s.seeds[name] = append(s.seeds[name], Seed{
		Name:    name,
		Version: version,
		Handler: handler,
	})
	return s
}

func (s *Seeder) AddSeedx(name string, version int, handler SeedHandlerx) *Seeder {
	if _, exists := s.seeds[name]; !exists {
		s.seeds[name] = make([]Seed, 0)
	}
	nHandler := func(ctx context.Context, tx zsql.Tx, seed Seed) error {
		txx, ok := tx.(zsql.Txx)
		if !ok {
			panic("AddSeedx used without zsql.DBx connection")
		}
		return handler(ctx, txx, seed)
	}
	s.seeds[name] = append(s.seeds[name], Seed{
		Name:    name,
		Version: version,
		Handler: nHandler,
	})
	return s
}

// Seed the given name to the specified version. Does nothing if the version is
// already satisfied. Can only seed forwards and not backwards
func (s *Seeder) SeedTo(name string, targetVersion int) (err error) {
	if err = s.init(); err != nil {
		return
	}
	// validate this key/version pair has been registered
	targetIndex := -1
	for i := range s.seeds[name] {
		if s.seeds[name][i].Version == targetVersion {
			targetIndex = i
			break
		}
	}
	if targetIndex < 0 {
		return fmt.Errorf("no seed registered for %s (%d)", name, targetVersion)
	}

	// determine the current seed version
	currentVersion, currentIndex, err := s.getCurrentVersion(name)
	if err != nil {
		return
	}

	// validate that the new version is higher than the current
	if currentVersion > targetVersion {
		return fmt.Errorf("target version of %d is not valid for current version %d.\n target version cannot be below current version", targetVersion, currentVersion)
	}
	// perform the seeds required to advance to the desired version
	seeds := s.seeds[name][currentIndex : targetIndex+1]
	ctx := context.Background()
	dbx, isDBx := s.db.(zsql.DBx)
	if isDBx {
		return zsql.WithBeginTxx(dbx, func(tx zsql.Txx) (err error) {
			if err = s.runSeeds(ctx, seeds, tx); err != nil {
				return
			}
			return s.insertSeed(tx, name, targetVersion)
		}, ctx, nil)
	} else {
		return zsql.WithBeginTx(s.db, func(tx zsql.Tx) (err error) {
			if err = s.runSeeds(ctx, seeds, tx); err != nil {
				return
			}
			return s.insertSeed(tx, name, targetVersion)
		}, ctx, nil)
	}
}

func (s *Seeder) insertSeed(tx zsql.Tx, name string, version int) error {
	q := fmt.Sprintf("INSERT INTO %s (name, version) VALUES (?, ?)", s.config.TableName)
	_, err := tx.Exec(q, name, version)
	return err
}

func (s *Seeder) init() (err error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if s.isInitialized {
		return
	}
	s.sortByVersion()
	q := fmt.Sprintf("SELECT * FROM %s LIMIT 1", s.config.TableName)
	r := s.db.QueryRow(q)
	if s.driver.IsNoTableErr(r.Err()) {
		q = fmt.Sprintf(seedSchema, s.config.TableName)
		if _, err = s.db.Exec(q); err != nil {
			return
		}
	}
	s.isInitialized = true
	return
}

// sort the registered seeds by version
func (s *Seeder) sortByVersion() {
	for key := range s.seeds {
		sort.Slice(s.seeds[key], func(a, b int) bool {
			return s.seeds[key][a].Version < s.seeds[key][b].Version
		})
	}
}

func (s *Seeder) runSeeds(ctx context.Context, seeds []Seed, tx zsql.Tx) (err error) {
	for _, seed := range seeds {
		if err = seed.Handler(ctx, tx, seed); err != nil {
			return
		}
	}
	return
}

func (s *Seeder) getCurrentVersion(name string) (version int, index int, err error) {
	q := fmt.Sprintf("SELECT version FROM %s WHERE name=? ORDER BY version DESC LIMIT 1", s.config.TableName)
	row := s.db.QueryRow(q, name)
	if err = row.Err(); err != nil {
		return
	}
	err = row.Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	q = fmt.Sprintf("SELECT count(*) as c FROM %s WHERE name=? ORDER BY version DESC LIMIT 1", s.config.TableName)
	row = s.db.QueryRow(q, name)
	if err = row.Err(); err != nil {
		return
	}
	err = row.Scan(&index)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func (s *Seeder) hasRegisteredSeed(name string, version int) bool {
	seeds, ok := s.seeds[name]
	if !ok {
		return false
	}
	for _, seed := range seeds {
		if seed.Version == version {
			return true
		}
	}
	return false
}
