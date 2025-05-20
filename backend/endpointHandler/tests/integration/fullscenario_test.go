//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Check first user-scenario:
// [Backend - Jira Connector - DB]:
// [Backend -> Jira Connector] GET api/v1/connector/projects?page=1&limit=9&search= -> получены 9 проектов из Jira
// [Backend -> JiraConnector -> DB] POST api/v1/connector/updateProject?project=KEY -> загружен один выбранный проект в базу данных
// [Backend -> DB] GET api/v1//projects-> загружен 1 проект с тем же именем и ключом
// [Backend -> DB] GET api/v1//projects/ID-> получение всей информации о проекте из базы данных
// [Backend -> DB] GET api/v1/analytics/status-distribution?key=KEY -> информация о соотношении отерытых\закрытых задач не пуста
// [Backend -> DB] DELETE api/v1/projects/ID-> удаление выбранного проекта по id из базы данных
func TestFullScenarioFirst(t *testing.T) {
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

	// check full info about project
	var prjStatus ProjectStats
	t.Run("GET /projects/{:id}", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/projects/%s", baseURL, loadProject[0].Id))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&prjStatus)
		require.NoError(t, err)
		require.NotEmpty(t, prjStatus)
	})

	// check analyticts status-distribution
	t.Run("GET /analytics/status-distribution", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/analytics/status-distribution?key=%s", baseURL, loadProject[0].Name))
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		t.Logf("Analytics response: %s", string(body))

		assert.Equal(t, 200, resp.StatusCode)
	})

	// delete project from DB
	t.Run("DELETE /deleteProject", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", baseURL, loadProject[0].Id), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
	})
}

// Check first user-scenario:
// [Backend - Jira Connector - DB]:
// [Backend -> Jira Connector] GET api/v1/connector/projects -> получены проекты из Jira
// [Backend -> JiraConnector -> DB] POST api/v1/connector/updateProject?project=KEY -> загружен один выбранный проект в базу данных
// [Backend -> JiraConnector -> DB] POST api/v1/connector/updateProject?project=KEY -> загружен второй выбранный проект в базу данных
// [Backend -> DB] GET api/v1//projects-> загружено 2 проекта с теми же именами
// [Backend -> DB] GET api/v1/compare/priority?key=KEY1,KEY2-> получение сравнительных данных по двум ранее загруженным задачам
func TestFullscenarioSecond(t *testing.T) {
	client := http.Client{Timeout: 20 * time.Second}

	baseURL := "http://backend:8000/api/v1"
	var selectedProjects ResponseProject

	// get all projects from jira
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

	// load two projects to DB
	t.Run("POST /connector/updateProject [two]", func(t *testing.T) {
		for _, key := range []string{p1.Key, p2.Key} {
			resp, err := client.Post(fmt.Sprintf("%s/connector/updateProject?project=%s", baseURL, key), "", nil)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, 200, resp.StatusCode)
		}
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
		require.Equal(t, len(loadProject), 2)
		require.Equal(t, p1.Name, loadProject[0].Name)
		require.Equal(t, p2.Name, loadProject[1].Name)
	})

	// compare two projects
	t.Run("GET /compare/priority", func(t *testing.T) {
		url := fmt.Sprintf("%s/compare/priority?key=%s,%s", baseURL, loadProject[0].Name, loadProject[1].Name)
		resp, err := client.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		var result PriorityStats
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)
		require.NotEmpty(t, result)
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
