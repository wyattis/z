package zmigrate

import (
	"context"
	_ "embed"
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

// Seed the given name to the specified version. Does nothing if the version is
// already satisfied. Can only seed forwards and not backwards
func (s *Seeder) SeedTo(name string, targetVersion int) (err error) {
	if err = s.init(); err != nil {
		return
	}
	// validate this key/version pair has been registered
	if !s.hasRegisteredSeed(name, targetVersion) {
		return fmt.Errorf("no seed registered for %s (%d)", name, targetVersion)
	}

	// determine the current seed version
	currentVersion, err := s.getCurrentVersion(name)
	if err != nil {
		return
	}
	// validate that the new version is higher than the current
	if currentVersion > targetVersion {
		return fmt.Errorf("target version of %d is not valid for current version %d.\n target version cannot be below current version", targetVersion, currentVersion)
	}
	// perform the seeds required to advance to the desired version
	return s.exec(s.seeds[name][currentVersion : targetVersion+1])
}

func (s *Seeder) init() (err error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if s.isInitialized {
		return
	}
	s.sortByVersion()
	fmt.Println(s.seeds)
	q := fmt.Sprintf("SELECT * FROM %s LIMIT 1", s.config.TableName)
	_, err = s.db.Exec(q)
	if s.driver.IsNoTableErr(err) {
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
			return s.seeds[key][a].Version > s.seeds[key][b].Version
		})
	}
}

func (s *Seeder) exec(seeds []Seed) error {
	ctx := context.Background()
	return zsql.WithBeginTx(s.db, func(tx zsql.Tx) (err error) {
		for _, seed := range seeds {
			if err = seed.Handler(ctx, tx, seed); err != nil {
				return
			}
		}
		return
	}, ctx, nil)
}

func (s *Seeder) getCurrentVersion(name string) (version int, err error) {
	q := fmt.Sprintf("SELECT version FROM %s WHERE name=?", name)
	row := s.db.QueryRowContext(context.Background(), q, name)
	if err = row.Err(); err != nil {
		return
	}
	err = row.Scan(&version)
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
