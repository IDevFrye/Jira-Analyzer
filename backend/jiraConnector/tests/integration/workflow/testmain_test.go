//go:build integration
// +build integration

package workflowintegrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jiraconnector/cmd/app"
	dbpusher "github.com/jiraconnector/internal/dbPusher"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Глобальные переменные для общего состояния
var (
	testApp    *app.JiraApp
	testConfig *config.Config
	mockJira   *httptest.Server
	testDB     *sql.DB
	testLogger *slog.Logger
)

const sqlCreatePath = "initDB/create.sql"

// TestMain настраивает все зависимости перед запуском тестов
func TestMain(m *testing.M) {
	// get config
	testConfig = prepareConfig()

	// set logger
	testLogger = logger.SetupLogger(testConfig.Env, testConfig.LogFile)

	// test databae
	postgresC, err := setupTestDB(testConfig, testLogger)
	if err != nil {
		testLogger.Error("Failed to start test DB", "error", err)
		os.Exit(1)
	}
	defer postgresC.Terminate(context.Background())

	// set test jira api
	mockJira = setupMockJira()
	defer mockJira.Close()

	testConfig.JiraCfg.Url = mockJira.URL

	// init my app
	testApp, err = app.NewApp(testConfig, testLogger)
	if err != nil {
		testLogger.Error("Failed to create app", "error", err)
		os.Exit(1)
	}
	defer testApp.Close()

	// start my app
	go func() {
		if err := testApp.Run(); err != nil {
			testLogger.Error("Server error", "error", err)
		}
	}()
	time.Sleep(500 * time.Millisecond)

	// run tests
	code := m.Run()
	os.Exit(code)
}

func setupTestDB(cfg *config.Config, log *slog.Logger) (testcontainers.Container, error) {
	ctx := context.Background()

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DBCfg.User,
			"POSTGRES_PASSWORD": cfg.DBCfg.Password,
			"POSTGRES_DB":       cfg.DBCfg.Name,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	testConfig.DBCfg.Host, _ = postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")
	testConfig.DBCfg.Port = port.Port()

	DB, err := dbpusher.NewDbPusher(cfg, log)
	if err != nil {
		panic(fmt.Errorf("failed to connect to test db: %w", err))
	}

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

	testDB = DB.Db()

	return postgresC, nil
}

func setupMockJira() *httptest.Server {
	projects := []structures.JiraProject{
		{Key: "TEST1", Name: "Test Project 1"},
		{Key: "TEST2", Name: "Test Project 2"},
	}

	issues := map[string][]structures.JiraIssue{
		"TEST1": {
			{Key: "TEST1-1", Fields: structures.Field{Summary: "Issue 1"}},
			{Key: "TEST1-2", Fields: structures.Field{Summary: "Issue 2"}},
		},
		"TEST2": {
			{Key: "TEST2-1", Fields: structures.Field{Summary: "Demo Issue"}},
		},
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testLogger.Info("Mock Jira request", "method", r.Method, "path", r.URL.Path)

		switch {
		case r.URL.Path == "/rest/api/2/project":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(projects)

		case r.URL.Path == "/rest/api/2/project/TEST1":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(projects[0])

		case strings.Contains(r.URL.Path, "/rest/api/2/search"):
			projectKey := strings.TrimPrefix(r.URL.Query().Get("jql"), "project=")
			if issues, ok := issues[projectKey]; ok {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(structures.JiraIssues{
					Issues: issues,
					Total:  len(issues),
				})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func prepareConfig() *config.Config {
	return &config.Config{
		Env:     "local",
		LogFile: "workflowintegrations.log",
		ServerCfg: config.ServerConfig{
			Port: ":8080",
		},
		JiraCfg: config.JiraConfig{
			Url:           "",
			ThreadCount:   2,
			IssueInOneReq: 50,
			MinSleep:      10,
			MaxSleep:      1000,
		},
		DBCfg: config.DBConfig{
			Host:     "0",
			Port:     "",
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
		},
	}
}
