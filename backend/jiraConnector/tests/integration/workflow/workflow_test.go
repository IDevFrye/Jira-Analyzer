//go:build integration
// +build integration

package workflowintegrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullIntegration(t *testing.T) {

	// /projects
	t.Run("Get projects list", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost%s/api/v1/connector/projects", testConfig.ServerCfg.Port))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Projects []structures.JiraProject `json:"projects"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Len(t, result.Projects, 2)
	})

	// /updateProject?project=NAME
	t.Run("Load project issues", func(t *testing.T) {
		projectKey := "TEST1"

		// Загружаем проект
		resp, err := http.Post(
			fmt.Sprintf("http://localhost%s/api/v1/connector/updateProject?project=%s", testConfig.ServerCfg.Port, projectKey),
			"application/json",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Проверяем БД
		var projectCount int
		err = testDB.QueryRow(
			"SELECT COUNT(*) FROM projects WHERE key = $1", projectKey,
		).Scan(&projectCount)
		require.NoError(t, err)
		assert.Equal(t, 1, projectCount)

		var issuesCount int
		err = testDB.QueryRow(
			"SELECT COUNT(*) FROM issue WHERE key LIKE $1", projectKey+"%",
		).Scan(&issuesCount)
		require.NoError(t, err)
		assert.Equal(t, 2, issuesCount)
	})

	// /updateProject?project=UNKNOWN
	t.Run("Error cases", func(t *testing.T) {
		resp, err := http.Post(
			fmt.Sprintf("http://localhost%s/api/v1/connector/updateProject?project=UNKNOWN", testConfig.ServerCfg.Port),
			"application/json",
			nil,
		)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
