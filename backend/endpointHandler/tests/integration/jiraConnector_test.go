//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// [Backend - Jira Connector] -> GET api/v1/connector/projects
func TestJiraConnectorProjects(t *testing.T) {
	client := http.Client{Timeout: 20 * time.Second}

	baseURL := getBackendBaseURL()
	var selectedProjects ResponseProject

	// get page=1 limit=9 projects from jira
	t.Run("GET /connector/projects", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/connector/projects?page=1&limit=9&search=", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&selectedProjects)
		require.NoError(t, err)
		require.Equal(t, len(selectedProjects.Projects), 9)
	})
}

// [Backend - Jira Connector] -> POST api/v1/connector/updateProject
func TestJiraConnectorUpdate(t *testing.T) {
	client := http.Client{Timeout: 20 * time.Second}

	baseURL := getBackendBaseURL()
	var selectedProjects ResponseProject

	// get page=1 limit=9 projects from jira
	t.Run("GET /connector/projects", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/connector/projects?page=1&limit=9&search=", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&selectedProjects)
		require.NoError(t, err)
		require.Equal(t, len(selectedProjects.Projects), 9)
	})

	p1 := selectedProjects.Projects[0]
	p2 := selectedProjects.Projects[1]
	p3 := selectedProjects.Projects[2]

	// load projects to DB to get at least one synced project
	t.Run("POST /connector/updateProject", func(t *testing.T) {
		for _, key := range []string{p1.Key, p2.Key, p3.Key} {
			resp, err := client.Post(fmt.Sprintf("%s/connector/updateProject?project=%s", baseURL, key), "", nil)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, 200, resp.StatusCode)
		}
	})

	// check that project in DB
	var loadProject []JiraProject
	var targetProject JiraProject
	t.Run("GET /projects", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/projects", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&loadProject)
		require.NoError(t, err)
		var found bool
		targetProject, found = pickFirstSyncedProject(loadProject, p1, p2, p3)
		require.True(t, found)
	})

	t.Run("DELETE /deleteProject", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", baseURL, targetProject.Id), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
	})
}
