//go:build integration
// +build integration

package jiraapiintegrations

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jiraconnector/internal/connector"
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
