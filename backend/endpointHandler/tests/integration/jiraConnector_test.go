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

	baseURL := "http://backend:8000/api/v1"
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

	baseURL := "http://backend:8000/api/v1"
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

	// load project to DB
	t.Run("POST /connector/updateProject", func(t *testing.T) {
		resp, err := client.Post(fmt.Sprintf("%s/connector/updateProject?project=%s", baseURL, p1.Key), "", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
	})

	// check that project in DB
	var loadProject []JiraProject
	t.Run("GET /projects", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/projects", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&loadProject)
		require.NoError(t, err)
		require.Equal(t, p1.Name, loadProject[0].Name)
	})

	t.Run("DELETE /deleteProject", func(t *testing.T) {
		for _, prj := range loadProject {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", baseURL, prj.Id), nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
		}
	})
}
