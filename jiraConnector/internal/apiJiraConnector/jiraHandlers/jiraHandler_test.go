package jirahandlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	jirahandlers "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers"
	"github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers/mocks"
	"github.com/jiraconnector/internal/structures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectsHandler(t *testing.T) {
	mockService := new(mocks.JiraServiceInterface)

	mockService.On("GetProjectsPage", "", 20, 1).
		Return(&structures.ResponseProject{
			Projects: []structures.JiraProject{{Id: "1", Name: "Test Project"}},
			PageInfo: structures.PageInfo{PageCount: 1, CurrentPage: 1, ProjectsCount: 1},
		}, nil)

	router := mux.NewRouter()
	handler := jirahandlers.NewHandler(mockService, router)

	req, _ := http.NewRequest("GET", "/projects", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response structures.ResponseProject
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Project", response.Projects[0].Name)

	mockService.AssertExpectations(t)
}

func TestUpdateProjectHandler_Success(t *testing.T) {
	mockService := new(mocks.JiraServiceInterface)

	mockService.On("UpdateProjects", "TestProject").
		Return([]structures.JiraIssue{{Key: "ISSUE-1"}}, nil)

	mockService.On("PushDataToDb", "TestProject", mock.Anything).
		Return(nil)

	router := mux.NewRouter()
	handler := jirahandlers.NewHandler(mockService, router)

	req, _ := http.NewRequest("POST", "/updateProject?project=TestProject", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := `{"TestProject":"updated"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())

	mockService.AssertExpectations(t)
}
