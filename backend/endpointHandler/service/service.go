package service

import (
	"encoding/json"
	"fmt"
	"github.com/endpointhandler/model"
	"github.com/endpointhandler/repository"
	"net/http"
)

func FetchJiraProjects() ([]model.Project, error) {
	resp, err := http.Get("http://localhost:8080/api/v1/connector/projects")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Projects []model.Project `json:"projects"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result.Projects, err
}

func UpdateJiraProject(project string) (map[string]string, error) {
	url := fmt.Sprintf("http://localhost:8080/api/v1/connector/updateProject?project=%s", project)
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
