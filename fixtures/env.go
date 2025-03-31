package fixtures

import (
	"async-api/config"
	"async-api/store"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

type TestEnv struct {
	Db     *sql.DB
	Config *config.Config
}

func NewTestEnv(t *testing.T) *TestEnv {
	_ = os.Setenv("ENV", string(config.Env_Test))
	cfg, err := config.New()
	require.NoError(t, err)

	db, err := store.NewPostgresDB(cfg)
	require.NoError(t, err)

	return &TestEnv{
		Db:     db,
		Config: cfg,
	}
}

func (te *TestEnv) SetupDb(t *testing.T) func(t *testing.T) {
	m, err := newMigrate(te)
	require.NoError(t, err)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}

	return te.TeardownDb
}

func (te *TestEnv) TeardownDb(t *testing.T) {
	_, err := te.Db.Exec(fmt.Sprintf("TRUNCATE TABLE %s",
		strings.Join([]string{"users", "reports", "refresh_tokens"}, ",")))
	require.NoError(t, err)
}

func newMigrate(te *TestEnv) (*migrate.Migrate, error) {
	m, err := migrate.New(
		fmt.Sprintf("file:///%s/migrations", te.Config.ProjectRoot),
		te.Config.DatabaseUrl())
	if err != nil {
		return nil, fmt.Errorf("could not create migration: %w", err)
	}
	return m, nil
}
