package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/endpointhandler/config"
	"github.com/endpointhandler/model"
	"github.com/endpointhandler/service"
	"github.com/stretchr/testify/assert"
)

func TestFetchJiraProjects_Success(t *testing.T) {
	expected := model.ProjectsResponse{
		Projects: []model.Project{
			{ID: "1", Key: "PROJ1", Name: "Project One", Self: "url1"},
		},
	}

	// Создаем мок-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects", r.URL.Path)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.Connector.BaseURL = server.URL

	resp, err := service.FetchJiraProjects(cfg)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestFetchJiraProjects_Failure(t *testing.T) {
	cfg := &config.Config{}
	cfg.Connector.BaseURL = "http://invalid.url" // несуществующий URL

	_, err := service.FetchJiraProjects(cfg)
	assert.Error(t, err)
}

func TestUpdateJiraProject_Success(t *testing.T) {
	responseMap := map[string]string{"status": "updated"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/updateProject", r.URL.Path)
		assert.Equal(t, "project=PROJ1", r.URL.RawQuery)
		json.NewEncoder(w).Encode(responseMap)
	}))
	defer server.Close()

	cfg := &config.Config{}
	cfg.Connector.BaseURL = server.URL

	result, err := service.UpdateJiraProject(cfg, "PROJ1")
	assert.NoError(t, err)
	assert.Equal(t, responseMap, result)
}

func TestUpdateJiraProject_Failure(t *testing.T) {
	cfg := &config.Config{}
	cfg.Connector.BaseURL = "http://invalid.url"

	_, err := service.UpdateJiraProject(cfg, "PROJ1")
	assert.Error(t, err)
}

func TestFetchAndStoreProjects(t *testing.T) {
	projectsResp := model.ProjectsResponse{
		Projects: []model.Project{
			{ID: "1", Key: "PROJ1", Name: "Project One", Self: "url1"},
		},
	}

	fetchServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(projectsResp)
	}))
	defer fetchServer.Close()

	updateServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}))
	defer updateServer.Close()

	cfg := &config.Config{}

	combinedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/projects" {
			json.NewEncoder(w).Encode(projectsResp)
		} else if r.URL.Path == "/updateProject" {
			json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
		} else {
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer combinedServer.Close()

	cfg.Connector.BaseURL = combinedServer.URL

	resp, err := service.FetchAndStoreProjects(cfg)
	assert.NoError(t, err)
	assert.Equal(t, projectsResp, resp)
}
