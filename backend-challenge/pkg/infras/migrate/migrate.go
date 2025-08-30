package migrate

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"

	// migrate
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	// migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrateTool interface {
	Migrate(source, connStr string) error
}

var onceIns sync.Once
var singleton MigrateTool
var mutex = &sync.Mutex{}

type migrateTool struct {
}

func New() MigrateTool {
	onceIns.Do(func() {
		singleton = &migrateTool{}
	})

	return singleton
}

func (t *migrateTool) Migrate(source string, connStr string) error {
	mutex.Lock()
	defer mutex.Unlock()

	m, err := migrate.New(source, connStr)
	if err != nil {
		return err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("get migrate version failed with err=%w", err)
	}

	if dirty {
		if err := m.Force(int(version) - 1); err != nil {
			return fmt.Errorf("failed to down version migration: err=%w", err)
		}
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}

	return nil
}
