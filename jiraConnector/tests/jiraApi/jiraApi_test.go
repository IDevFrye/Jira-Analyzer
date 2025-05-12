package jiraapiintegrations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jiraconnector/internal/connector"
	"github.com/jiraconnector/internal/structures"
	"github.com/jiraconnector/pkg/config"
	"github.com/jiraconnector/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllProjects(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	cfg := config.Config{
		JiraCfg: config.JiraConfig{
			Url:           ts.URL,
			MinSleep:      100,
			MaxSleep:      1000,
			ThreadCount:   2,
			IssueInOneReq: 50,
		},
	}
	log := logger.SetupLogger("test", "")

	conn := connector.NewJiraConnector(&cfg, log)

	projects, err := conn.GetAllProjects()
	require.NoError(t, err)
	assert.Len(t, projects, 3)
	assert.Equal(t, "TEST1", projects[0].Key)
}

func TestGetProjectsPage(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	cfg := config.Config{JiraCfg: config.JiraConfig{Url: ts.URL}}
	log := logger.SetupLogger("test", "")

	conn := connector.NewJiraConnector(&cfg, log)

	t.Run("First page", func(t *testing.T) {
		resp, err := conn.GetProjectsPage("", 2, 1)
		require.NoError(t, err)
		assert.Len(t, resp.Projects, 2)
		assert.Equal(t, 2, resp.PageInfo.PageCount)
	})

	t.Run("Search filtered", func(t *testing.T) {
		resp, err := conn.GetProjectsPage("test", 10, 1)
		require.NoError(t, err)
		assert.Len(t, resp.Projects, 2) // Only TEST1 and TEST2 match
	})

	t.Run("Empty page", func(t *testing.T) {
		resp, err := conn.GetProjectsPage("", 2, 3)
		require.NoError(t, err)
		assert.Empty(t, resp.Projects)
	})
}

func TestGetProjectIssues(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	cfg := config.Config{
		JiraCfg: config.JiraConfig{
			Url:           ts.URL,
			ThreadCount:   3,
			IssueInOneReq: 50,
			MinSleep:      100,
			MaxSleep:      1000,
		},
	}
	log := logger.SetupLogger("debug", "jiraApiintegrations.log")

	conn := connector.NewJiraConnector(&cfg, log)

	t.Run("Single thread", func(t *testing.T) {
		cfg.JiraCfg.ThreadCount = 1
		_, err := conn.GetProjectIssues("TEST1")
		require.NoError(t, err)
		//assert.Len(t, issues, 50) // Only first batch
	})

	t.Run("Multiple threads", func(t *testing.T) {
		cfg.JiraCfg.ThreadCount = 3
		_, err := conn.GetProjectIssues("TEST1")
		require.NoError(t, err)
		//assert.Len(t, issues, 150) // All issues from 3 threads
	})

	t.Run("Empty project", func(t *testing.T) {
		// Modify test server handler for this case
		ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/rest/api/2/search") {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(structures.JiraIssues{Total: 0})
				return
			}
			w.WriteHeader(http.StatusNotFound)
		})

		issues, err := conn.GetProjectIssues("EMPTY")
		require.NoError(t, err)
		assert.Empty(t, issues)
	})
}

func TestErrorHandling(t *testing.T) {
	log := logger.SetupLogger("debug", "jiraApiintegrations.log")
	t.Run("Jira unavailable", func(t *testing.T) {
		cfg := config.Config{JiraCfg: config.JiraConfig{Url: "http://invalid"}}
		conn := connector.NewJiraConnector(&cfg, log)

		_, err := conn.GetAllProjects()
		assert.Error(t, err)
	})

	t.Run("Rate limiting", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer ts.Close()

		cfg := config.Config{
			JiraCfg: config.JiraConfig{
				Url:      ts.URL,
				MinSleep: 10,
				MaxSleep: 50,
			},
		}
		conn := connector.NewJiraConnector(&cfg, log)

		_, err := conn.GetAllProjects()
		assert.Error(t, err)
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer ts.Close()

		cfg := config.Config{JiraCfg: config.JiraConfig{Url: ts.URL}}
		conn := connector.NewJiraConnector(&cfg, log)

		_, err := conn.GetAllProjects()
		assert.Error(t, err)
	})
}

func TestRetryRequest(t *testing.T) {
	var attempt int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`[{"key":"TEST","name":"Test Project"}]`))
	}))
	defer ts.Close()

	cfg := config.Config{
		JiraCfg: config.JiraConfig{
			Url:      ts.URL,
			MinSleep: 10,
			MaxSleep: 100,
		},
	}
	log := logger.SetupLogger("test", "")
	conn := connector.NewJiraConnector(&cfg, log)

	projects, err := conn.GetAllProjects()
	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, 3, attempt)
}

func TestRetryRequestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := config.Config{
		JiraCfg: config.JiraConfig{
			Url:      ts.URL,
			MinSleep: 10,
			MaxSleep: 20,
		},
	}
	log := logger.SetupLogger("test", "")
	conn := connector.NewJiraConnector(&cfg, log)

	_, err := conn.GetAllProjects()
	assert.Error(t, err)
}

func TestIntegration(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	cfg := config.Config{
		JiraCfg: config.JiraConfig{
			Url:           ts.URL,
			ThreadCount:   2,
			IssueInOneReq: 50,
			MinSleep:      10,
			MaxSleep:      100,
		},
	}
	log := logger.SetupLogger("test", "")
	conn := connector.NewJiraConnector(&cfg, log)

	t.Run("Full flow", func(t *testing.T) {
		// 1. Get all projects
		projects, err := conn.GetAllProjects()
		require.NoError(t, err)
		require.Len(t, projects, 3)

		// 2. Get paginated projects
		page, err := conn.GetProjectsPage("test", 1, 1)
		require.NoError(t, err)
		assert.Len(t, page.Projects, 1)

		// 3. Get issues for project
		issues, err := conn.GetProjectIssues("TEST1")
		require.NoError(t, err)
		assert.Len(t, issues, 0)
	})
}

func TestPerformance(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	cfg := config.Config{
		JiraCfg: config.JiraConfig{
			Url:           ts.URL,
			ThreadCount:   5,
			IssueInOneReq: 100,
		},
	}
	conn := connector.NewJiraConnector(&cfg, logger.SetupLogger("debug", "jiraApiintegrations.log"))

	start := time.Now()
	_, err := conn.GetProjectIssues("PERF")
	assert.NoError(t, err)
	assert.Less(t, time.Since(start).Seconds(), 1.0, "Should complete in under 1 second")
}
