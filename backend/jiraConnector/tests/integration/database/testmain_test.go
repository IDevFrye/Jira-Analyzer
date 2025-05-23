//go:build integration
// +build integration

package dbintegrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	dbpusher "github.com/jiraconnector/internal/dbPusher"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var DB *dbpusher.DbPusher

const sqlCreatePath = "initDB/create.sql"

func TestMain(m *testing.M) {
	ctx := context.Background()

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start postgres container: %w", err))
	}
	defer postgresC.Terminate(ctx)

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	cfg := config.Config{
		Env:     "debag",
		LogFile: "integrationDB.log",
		DBCfg: config.DBConfig{
			Host:     host,
			Port:     port.Port(),
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
		},
	}

	log := logger.SetupLogger(cfg.Env, cfg.LogFile)

	DB, err = dbpusher.NewDbPusher(&cfg, log)
	if err != nil {
		panic(fmt.Errorf("failed to connect to test db: %w", err))
	}
	defer DB.Close()

	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	sqlPath := filepath.Join(projectRoot, "build", sqlCreatePath)
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		panic(fmt.Errorf("failed to read create.sql: %w", err))
	}

	_, err = DB.Db().Exec(string(sqlBytes))
	if err != nil {
		panic(fmt.Errorf("failed to execute schema.sql: %w", err))
	}

	os.Exit(m.Run())
}

func resetTestDB(t *testing.T) {
	_, err := DB.Db().Exec("DROP SCHEMA public CASCADE")
	require.NoError(t, err)

	_, err = DB.Db().Exec("CREATE SCHEMA public")
	require.NoError(t, err)

	// Повторно применяем миграции
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	sqlPath := filepath.Join(projectRoot, "build", sqlCreatePath)
	sqlBytes, err := os.ReadFile(sqlPath)
	require.NoError(t, err)

	_, err = DB.Db().Exec(string(sqlBytes))
	require.NoError(t, err)
}
