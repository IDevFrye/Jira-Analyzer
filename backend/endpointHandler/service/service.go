package service

import (
	"encoding/json"
	"fmt"
	"github.com/endpointhandler/config"
	"net/http"

	"github.com/endpointhandler/model"
	"github.com/endpointhandler/repository"
)

func FetchJiraProjects(cfg *config.Config) (model.ProjectsResponse, error) {
	url := fmt.Sprintf("%s/projects", cfg.Connector.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return model.ProjectsResponse{}, err
	}
	defer resp.Body.Close()

	var result model.ProjectsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func UpdateJiraProject(cfg *config.Config, project string) (map[string]string, error) {
	url := fmt.Sprintf("%s/updateProject?project=%s", cfg.Connector.BaseURL, project)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func GetAllProjects() ([]model.Project, error) {
	return repository.GetAllProjects()
}

func GetProjectStats(id int) (model.ProjectStats, error) {
	return repository.GetStats(id)
}

func DeleteProject(id int) error {
	return repository.DeleteProject(id)
}

func FetchAndStoreProjects(cfg *config.Config) (model.ProjectsResponse, error) {
	projectsResp, err := FetchJiraProjects(cfg)
	if err != nil {
		return model.ProjectsResponse{}, err
	}

	for _, p := range projectsResp.Projects {
		_, err := UpdateJiraProject(cfg, p.Key)
		if err != nil {
			return model.ProjectsResponse{}, fmt.Errorf("updateProject failed for key=%s: %w", p.Key, err)
		}
	}

	return projectsResp, nil
}
