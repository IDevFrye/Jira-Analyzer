package connector_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	configreader "github.com/jiraconnector/internal/configReader"
	"github.com/jiraconnector/internal/connector"
	"github.com/jiraconnector/internal/connector/mocks"
	"github.com/jiraconnector/internal/structures"

	"github.com/stretchr/testify/assert"
)

func TestGetAllProjects(t *testing.T) {
	testProjects := []structures.JiraProject{
		{Id: "1", Name: "Project_1", Key: "Key_1", Self: "url_1"},
		{Id: "2", Name: "Project_2", Key: "Key_2", Self: "url_2"},
	}

	//mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/project", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testProjects)
	}))
	defer server.Close()

	//test config
	config := configreader.Config{
		JiraCfg: configreader.JiraConfig{
			Url:           server.URL,
			ThreadCount:   1,
			IssueInOneReq: 100,
			MinSleep:      50,
			MaxSleep:      500,
		},
	}
	conn := connector.NewJiraConnector(config)

	//test method
	projects, err := conn.GetAllProjects()
	assert.NoError(t, err)
	assert.Equal(t, len(testProjects), len(projects))
	assert.Equal(t, testProjects[0].Name, projects[0].Name)
}

func TestGetProjectsPage(t *testing.T) {
	mockConnector := new(mocks.JiraConnectorInterface)

	mockConnector.On("GetAllProjects").Return([]structures.JiraProject{
		{Id: "1", Name: "Project_1", Key: "Key_1", Self: "url_1"},
		{Id: "2", Name: "Project_2", Key: "Key_2", Self: "url_2"},
		{Id: "3", Name: "Project_3", Key: "Key_3", Self: "url_3"},
	}, nil)

	mockConnector.On("GetProjectsPage", "Project", 2, 2).Return(&structures.ResponseProject{
		Projects: []structures.JiraProject{
			{Id: "2", Name: "Project_2", Key: "Key_2", Self: "url_2"},
			{Id: "3", Name: "Project_3", Key: "Key_3", Self: "url_3"},
		},
		PageInfo: structures.PageInfo{
			PageCount:     2,
			CurrentPage:   2,
			ProjectsCount: 3,
		},
	}, nil)

	result, err := mockConnector.GetProjectsPage("Project", 2, 2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Projects))
	assert.Equal(t, "Project_2", result.Projects[0].Name)
	assert.Equal(t, 2, result.PageInfo.CurrentPage)

	mockConnector.AssertCalled(t, "GetProjectsPage", "Project", 2, 2)
}

func TestGetProjectIssues(t *testing.T) {

	testIssues := structures.JiraIssues{
		StartAt:    1,
		MaxResults: 3,
		Total:      3,
		Issues: []structures.JiraIssue{
			{Id: "1", Key: "TEST-1", Fields: structures.Field{}},
			{Id: "2", Key: "TEST-2", Fields: structures.Field{}},
			{Id: "3", Key: "TEST-3", Fields: structures.Field{}},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testIssues)
	}))
	defer server.Close()

	config := configreader.Config{
		JiraCfg: configreader.JiraConfig{
			Url:           server.URL,
			ThreadCount:   1,
			IssueInOneReq: 50,
			MinSleep:      50,
			MaxSleep:      500,
		},
	}
	conn := connector.NewJiraConnector(config)

	issues, err := conn.GetProjectIssues("TEST")
	assert.NoError(t, err)
	assert.Equal(t, len(testIssues.Issues), len(issues))
}

func TestGetAllProjectsInvalidJSON(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	config := configreader.Config{
		JiraCfg: configreader.JiraConfig{
			Url: server.URL,
		},
	}
	conn := connector.NewJiraConnector(config)

	projects, err := conn.GetAllProjects()
	assert.Error(t, err)
	assert.Nil(t, projects)
}
